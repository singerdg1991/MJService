package types

/*
 * @apiDefine: UploadMetadata
 */
type UploadMetadata struct {
	FileName string `json:"fileName" openapi:"example:file.jpg"`
	FileSize int64  `json:"fileSize" openapi:"example:547207"`
	Path     string `json:"path,omitempty" openapi:"example:/uploads/staff"`
}
