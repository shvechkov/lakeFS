package esti

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/api/apiutil"
	gtwerrors "github.com/treeverse/lakefs/pkg/gateway/errors"
	"github.com/treeverse/lakefs/pkg/testutil"
)

type GetCredentials = func(id, secret, token string) *credentials.Credentials

const (
	numUploads           = 100
	randomDataPathLength = 1020
	gatewayTestPrefix    = "main/data/"
)

func newMinioClient(t *testing.T, getCredentials GetCredentials) *minio.Client {
	t.Helper()
	accessKeyID := viper.GetString("access_key_id")
	secretAccessKey := viper.GetString("secret_access_key")
	endpoint := viper.GetString("s3_endpoint")
	endpointSecure := viper.GetBool("s3_endpoint_secure")
	creds := getCredentials(accessKeyID, secretAccessKey, "")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: endpointSecure,
	})
	if err != nil {
		t.Fatalf("minio.New: %s", err)
	}
	return client
}

func TestS3UploadAndDownload(t *testing.T) {
	const parallelism = 10

	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)

	sigs := []struct {
		Name           string
		GetCredentials GetCredentials
	}{
		{Name: "V2", GetCredentials: credentials.NewStaticV2},
		{Name: "V4", GetCredentials: credentials.NewStaticV4},
	}

	for _, sig := range sigs {
		t.Run("Sig"+sig.Name, func(t *testing.T) {
			// Use the same sequence of path names to test each sig.
			r := rand.New(rand.NewSource(17))

			type Object struct {
				Path, Content string
			}

			var (
				wg      sync.WaitGroup
				objects = make(chan Object, parallelism*2)
			)

			for i := 0; i < parallelism; i++ {
				client := newMinioClient(t, sig.GetCredentials)

				wg.Add(1)
				go func() {
					for o := range objects {
						_, err := client.PutObject(
							ctx, repo, o.Path, strings.NewReader(o.Content), int64(len(o.Content)), minio.PutObjectOptions{})
						if err != nil {
							t.Errorf("minio.Client.PutObject(%s): %s", o.Path, err)
						}

						download, err := client.GetObject(
							ctx, repo, o.Path, minio.GetObjectOptions{})
						if err != nil {
							t.Errorf("minio.Client.GetObject(%s): %s", o.Path, err)
						}
						contents := bytes.NewBuffer(nil)
						_, err = io.Copy(contents, download)
						if err != nil {
							t.Errorf("download %s: %s", o.Path, err)
						}
						if strings.Compare(contents.String(), o.Content) != 0 {
							t.Errorf(
								"Downloaded bytes %v from uploaded bytes %v", contents.Bytes(), o.Content)
						}
					}
					wg.Done()
				}()
			}

			for i := 0; i < numUploads; i++ {
				objects <- Object{
					Content: testutil.RandomString(r, randomDataContentLength),
					// lakeFS supports _any_ path, even if its
					// byte sequence is not legal UTF-8 string.
					Path: gatewayTestPrefix + testutil.RandomString(r, randomDataPathLength-len(gatewayTestPrefix)),
				}
			}
			close(objects)
			wg.Wait()
		})
	}
}

func verifyObjectInfo(t *testing.T, got minio.ObjectInfo, expectedSize int) {
	if got.Err != nil {
		t.Errorf("%s: %s", got.Key, got.Err)
	}
	if got.Size != int64(expectedSize) {
		t.Errorf("Got size %d != expected size %d", got.Size, int64(expectedSize))
	}
}

func TestS3ReadObject(t *testing.T) {
	const (
		contents = "the quick brown fox jumps over the lazy dog"
		goodPath = "main/exists"
		badPath  = "main/does/not/exist"
	)

	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)

	// Upload an object
	minioClient := newMinioClient(t, credentials.NewStaticV2)

	_, err := minioClient.PutObject(ctx, repo, goodPath, strings.NewReader(contents), int64(len(contents)), minio.PutObjectOptions{})
	if err != nil {
		t.Errorf("PutObject(%s, %s): %s", repo, goodPath, err)
	}

	t.Run("get_exists", func(t *testing.T) {
		res, err := minioClient.GetObject(ctx, repo, goodPath, minio.GetObjectOptions{})
		if err != nil {
			t.Errorf("GetObject(%s, %s): %s", repo, goodPath, err)
		}
		defer func() { _ = res.Close() }()
		info, err := res.Stat()
		if err != nil {
			t.Errorf("GetObject(%s, %s) get: %s", repo, goodPath, err)
		}
		verifyObjectInfo(t, info, len(contents))
		got, err := io.ReadAll(res)
		if err != nil {
			t.Errorf("Read: %s", err)
		}
		if string(got) != contents {
			t.Errorf("Got contents \"%s\" but expected \"%s\"", string(got), contents)
		}
	})

	t.Run("get_exists_presigned", func(t *testing.T) {
		// using presigned URL, so we can check the headers
		// We expect the Content-Length header to be set
		preSignedURL, err := minioClient.Presign(ctx, http.MethodGet, repo, goodPath, time.Second*60, url.Values{})
		if err != nil {
			t.Fatalf("minioClientPresign(%s, %s): %s", repo, goodPath, err)
		}

		resp, err := http.Get(preSignedURL.String())
		if err != nil {
			t.Fatalf("http.Get %s: %s", preSignedURL.String(), err)
		}
		defer func() { _ = resp.Body.Close() }()
		contentLength := resp.Header.Get("Content-Length")
		if contentLength == "" {
			t.Errorf("Expected Content-Length header to be set")
		}
	})

	t.Run("get_no_physical_object", func(t *testing.T) {
		blockStoreType := viper.GetString("blockstore_type")
		if blockStoreType != "s3" {
			t.Skip("Skipping test - blockstore type is not s3")
		}
		awsAccessKeyID := viper.GetString("aws_access_key_id")
		awsSecretAccessKey := viper.GetString("aws_secret_access_key")
		if awsAccessKeyID == "" || awsSecretAccessKey == "" {
			t.Skip("Skipping test - AWS credentials not set")
		}

		// prepare entry with no physical object under it
		const (
			objPath   = "no-physical-object"
			objS3Path = "main/" + objPath
		)
		_, err = minioClient.PutObject(ctx, repo, objS3Path, strings.NewReader(""), int64(0), minio.PutObjectOptions{})
		if err != nil {
			t.Errorf("PutObject(%s, %s): %s", repo, objS3Path, err)
		}
		statResp, err := client.StatObjectWithResponse(ctx, repo, "main", &apigen.StatObjectParams{
			Path:    objPath,
			Presign: apiutil.Ptr(false),
		})
		if err != nil {
			t.Fatalf("StatObject(%s, %s): %s", repo, objPath, err)
		}
		if statResp.JSON200 == nil {
			t.Fatalf("StatObject(%s, %s): got nil response, status: %s", repo, objPath, statResp.Status())
		}
		physicalAddress := statResp.JSON200.PhysicalAddress

		// setup s3 client and delete the physical object
		s3BlockstoreClient, err := minio.New("s3.amazonaws.com", &minio.Options{
			Creds:  credentials.NewStaticV4(awsAccessKeyID, awsSecretAccessKey, ""),
			Secure: true,
		})
		if err != nil {
			t.Fatalf("failed to create s3 blockstore client: %s", err)
		}

		// access the underlying storage directly to delete the object
		s3PhysicalURL, err := url.Parse(physicalAddress)
		if err != nil {
			t.Fatalf("failed to parse physical address %s, %v", physicalAddress, err)
		}
		if s3PhysicalURL.Scheme != "s3" {
			t.Fatalf("physical address %s is not an s3 address", physicalAddress)
		}
		s3PhysicalKey := strings.TrimLeft(s3PhysicalURL.Path, "/")
		err = s3BlockstoreClient.RemoveObject(ctx, s3PhysicalURL.Host, s3PhysicalKey, minio.RemoveObjectOptions{})
		if err != nil {
			t.Fatalf("RemoveObject(%s, %s): %s", repo, objPath, err)
		}

		// try to read the object - should fail
		const s3ObjPath = "main/" + objPath
		res, err := minioClient.GetObject(ctx, repo, s3ObjPath, minio.GetObjectOptions{})
		if err != nil {
			t.Fatalf("GetObject(%s, %s): %s", repo, s3ObjPath, err)
		}
		defer func() { _ = res.Close() }()
		got, err := io.ReadAll(res)
		if err == nil {
			t.Fatalf("Successfully read \"%s\" from nonexistent path %s", got, s3ObjPath)
		}
		s3ErrorResponse := minio.ToErrorResponse(err)
		expectedErrorCode := gtwerrors.Codes[gtwerrors.ErrNoSuchVersion].Code
		if s3ErrorResponse.Code != expectedErrorCode {
			t.Errorf("Got %+v [%d] on reading when expecting code %s", s3ErrorResponse, s3ErrorResponse.StatusCode, expectedErrorCode)
		}
	})

	t.Run("get_not_exists", func(t *testing.T) {
		res, err := minioClient.GetObject(ctx, repo, badPath, minio.GetObjectOptions{})
		if err != nil {
			t.Errorf("GetObject(%s, %s): %s", repo, badPath, err)
		}
		defer func() { _ = res.Close() }()
		got, err := io.ReadAll(res)
		if err == nil {
			t.Errorf("Successfully read \"%s\" from nonexistent path %s", got, badPath)
		}
		s3ErrorResponse := minio.ToErrorResponse(err)
		if s3ErrorResponse.StatusCode != 404 {
			t.Errorf("Got %+v [%d] on reading when expecting Not Found [404]",
				s3ErrorResponse, s3ErrorResponse.StatusCode)
		}
	})

	t.Run("head_exists", func(t *testing.T) {
		info, err := minioClient.StatObject(ctx, repo, goodPath, minio.StatObjectOptions{})
		if err != nil {
			t.Errorf("StatObject(%s, %s): %s", repo, goodPath, err)
		}
		verifyObjectInfo(t, info, len(contents))
	})

	t.Run("head_not_exists", func(t *testing.T) {
		info, err := minioClient.StatObject(ctx, repo, badPath, minio.StatObjectOptions{})
		if err == nil {
			t.Errorf("StatObject(%s, %s): expected an error but got %+v", repo, badPath, info)
		}

		s3ErrorResponse := minio.ToErrorResponse(err)
		if s3ErrorResponse.StatusCode != 404 {
			t.Errorf("Got %+v [%d] on reading when expecting Not Found [404]",
				s3ErrorResponse, s3ErrorResponse.StatusCode)
		}
	})
}

func TestS3HeadBucket(t *testing.T) {
	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)

	badRepo := repo + "-nonexistent"

	// Upload an object
	client := newMinioClient(t, credentials.NewStaticV2)

	t.Run("existing", func(t *testing.T) {
		ok, err := client.BucketExists(ctx, repo)
		if err != nil {
			t.Errorf("BucketExists(%s) failed: %s", repo, err)
		}
		if !ok {
			t.Errorf("Got that good bucket %s does not exist", repo)
		}
	})

	t.Run("not existing", func(t *testing.T) {
		ok, err := client.BucketExists(ctx, badRepo)
		if err != nil {
			t.Errorf("BucketExists(%s) failed: %s", badRepo, err)
		}
		if ok {
			t.Errorf("Got that bad bucket %s exists", badRepo)
		}
	})
}

func TestS3CopyObject(t *testing.T) {
	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)

	// additional repository for copy between repos
	const destRepoName = "tests3copyobjectdest"
	destRepo := createRepositoryByName(ctx, t, destRepoName)
	defer deleteRepositoryIfAskedTo(ctx, destRepoName)

	// content
	r := rand.New(rand.NewSource(17))
	objContent := testutil.RandomString(r, randomDataContentLength)
	srcPath := gatewayTestPrefix + "source-file"
	destPath := gatewayTestPrefix + "dest-file"

	// upload data
	s3lakefsClient := newMinioClient(t, credentials.NewStaticV2)
	_, err := s3lakefsClient.PutObject(ctx, repo, srcPath, strings.NewReader(objContent), int64(len(objContent)),
		minio.PutObjectOptions{})
	require.NoError(t, err)

	t.Run("same_branch", func(t *testing.T) {
		// copy the object to the same repository
		_, err = s3lakefsClient.CopyObject(ctx,
			minio.CopyDestOptions{
				Bucket: repo,
				Object: destPath,
			},
			minio.CopySrcOptions{
				Bucket: repo,
				Object: srcPath,
			})
		if err != nil {
			t.Fatalf("minio.Client.CopyObject from(%s) to(%s): %s", srcPath, destPath, err)
		}

		download, err := s3lakefsClient.GetObject(ctx, repo, destPath, minio.GetObjectOptions{})
		if err != nil {
			t.Fatalf("minio.Client.GetObject(%s): %s", destPath, err)
		}
		defer func() { _ = download.Close() }()

		// compere files content
		var content bytes.Buffer
		_, err = io.Copy(&content, download)
		if err != nil {
			t.Fatalf("Download '%s' failed: %s", destPath, err)
		}
		require.Equal(t, objContent, content.String())

		resp, err := client.StatObjectWithResponse(ctx, repo, mainBranch, &apigen.StatObjectParams{Path: "data/source-file"})
		require.NoError(t, err)
		require.NotNil(t, resp.JSON200)

		resp, err = client.StatObjectWithResponse(ctx, repo, mainBranch, &apigen.StatObjectParams{Path: "data/dest-file"})
		require.NoError(t, err)
		require.NotNil(t, resp.JSON200)

		// assert that the physical addresses of the objects are the same
		sourceObjectStats := resp.JSON200
		destObjectStats := resp.JSON200
		require.Equal(t, sourceObjectStats.PhysicalAddress, destObjectStats.PhysicalAddress, "source and dest physical address should match")
	})

	t.Run("different_repo", func(t *testing.T) {
		// copy the object to different repository. should create another version of the file
		_, err := s3lakefsClient.CopyObject(ctx,
			minio.CopyDestOptions{
				Bucket: destRepo,
				Object: destPath,
			},
			minio.CopySrcOptions{
				Bucket: repo,
				Object: srcPath,
			})
		if err != nil {
			t.Fatalf("minio.Client.CopyObjectFrom(%s)To(%s): %s", srcPath, destPath, err)
		}

		download, err := s3lakefsClient.GetObject(ctx, destRepo, destPath, minio.GetObjectOptions{})
		if err != nil {
			t.Fatalf("minio.Client.GetObject(%s): %s", destPath, err)
		}
		defer func() { _ = download.Close() }()
		var contents bytes.Buffer
		_, err = io.Copy(&contents, download)
		require.NoError(t, err)

		// compere files content
		require.Equal(t, contents.String(), objContent)

		resp, err := client.StatObjectWithResponse(ctx, repo, mainBranch, &apigen.StatObjectParams{Path: "data/source-file"})
		require.NoError(t, err)
		require.NotNil(t, resp.JSON200)
		sourceObjectStats := resp.JSON200

		resp, err = client.StatObjectWithResponse(ctx, destRepo, mainBranch, &apigen.StatObjectParams{Path: "data/dest-file"})
		if err != nil {
			t.Fatalf("client.StatObject(%s): %s", destPath, err)
		}
		require.NoError(t, err)
		require.NotNil(t, resp.JSON200)
		destObjectStats := resp.JSON200

		// assert that the physical addresses of the objects are not the same
		require.NotEqual(t, sourceObjectStats.PhysicalAddress, destObjectStats.PhysicalAddress)
	})
}
