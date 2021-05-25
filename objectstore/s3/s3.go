package s3

import (
	"io"

	"github.com/adityak368/swissknife/logger/v2"
	"github.com/adityak368/swissknife/objectstore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3ObjectStore defines object store for aws s3
type S3ObjectStore struct {
	s3Uploader   *s3manager.Uploader
	s3Downloader *s3manager.Downloader
	s3Session    *session.Session
	s3Client     *s3.S3
}

func (store *S3ObjectStore) initStore() {
	// Create a single AWS session (we can re use this if we're uploading many files)
	s3Session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Store the session
	store.s3Session = s3Session

	// Create an uploader with the session and default options
	store.s3Uploader = s3manager.NewUploader(s3Session)

	// Create a downloader with the session and default options
	store.s3Downloader = s3manager.NewDownloader(s3Session)

	// Create a s3 client for object deletion
	store.s3Client = s3.New(s3Session)
	logger.Info().Msg("Initialized S3")
}

// AddFile will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (store *S3ObjectStore) AddFile(bucket, fileName string, src io.Reader) (*objectstore.UploadResult, error) {

	// Upload the file to S3.
	result, err := store.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		ACL:    aws.String("private"),
		Body:   src,
	})

	return &objectstore.UploadResult{
		Location: result.Location,
		UploadID: result.UploadID,
	}, err
}

// AddImage will upload a single image to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (store *S3ObjectStore) AddImage(bucket, fileName string, src io.Reader) (*objectstore.UploadResult, error) {

	// Upload the file to S3.
	result, err := store.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		ACL:         aws.String("private"),
		Body:        src,
		ContentType: aws.String("image/*"),
	})

	return &objectstore.UploadResult{
		Location: result.Location,
		UploadID: result.UploadID,
	}, err
}

// DeleteFile will delete a file from S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (store *S3ObjectStore) DeleteFile(bucket, fileName string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	_, err := store.s3Client.DeleteObject(input)
	return err
}

// DownloadFile will download a file from S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (store *S3ObjectStore) DownloadFile(bucket, fileName string, file io.WriterAt) (int64, error) {

	return store.s3Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
}

// New Creates a new S3 Store
func New() *S3ObjectStore {
	s3Store := new(S3ObjectStore)
	s3Store.initStore()
	return s3Store
}

var s3Store *S3ObjectStore

// Store is a singleton for plug and play usage
func Store() objectstore.ObjectStore {
	if s3Store == nil {
		s3Store = New()
	}
	return s3Store
}
