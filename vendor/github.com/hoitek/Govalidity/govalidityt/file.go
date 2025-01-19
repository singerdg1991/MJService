package govalidityt

import "mime/multipart"

type File struct {
	File   *multipart.File       `json:"file"`
	Header *multipart.FileHeader `json:"handler"`
}
