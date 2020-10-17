package objectstore

import "io"

// UploadResult defines the upload result
type UploadResult struct {
	Location string
	UploadID string
}

// ObjectStore defines the object store interface
type ObjectStore interface {
	AddFile(bucket, fileName string, src io.Reader) (*UploadResult, error)
	AddImage(bucket, fileName string, src io.Reader) (*UploadResult, error)
	DeleteFile(bucket, fileName string) error
	DownloadFile(bucket, fileName string, file io.WriterAt) (int64, error)
}
