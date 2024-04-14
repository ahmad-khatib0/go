package types

// Media handling constants
const (
	// UploadStarted indicates that the upload has started but not finished yet.
	UploadStarted = iota
	// UploadCompleted indicates that the upload has completed successfully.
	UploadCompleted
	// UploadFailed indicates that the upload has failed.
	UploadFailed
	// UploadDeleted indicates that the upload is no longer needed and can be deleted.
	UploadDeleted
)

// FileDef is a stored record of a file upload
type FileDef struct {
	ObjHeader `bson:",inline"`
	// Status of upload
	Status int
	// User who created the file
	User string
	// Type of the file.
	MimeType string
	// Size of the file in bytes.
	Size int64
	// Internal file location, i.e. path on disk or an S3 blob address.
	Location string
}
