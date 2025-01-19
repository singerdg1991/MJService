package utils

import (
	"fmt"
	govalidity "github.com/hoitek/Govalidity"
	"github.com/hoitek/Govalidity/govalidityt"
)

// ConvertNumberToHumanReadableSize converts a number to a human readable size
func ConvertNumberToHumanReadableSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%d KB", size/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%d MB", size/1024/1024)
	} else if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%d GB", size/1024/1024/1024)
	} else {
		return fmt.Sprintf("%d TB", size/1024/1024/1024/1024)
	}
}

// ValidateUploadFilesSize validates the size of the uploaded files
func ValidateUploadFilesSize(label string, files []*govalidityt.File, acceptedSizeInBytes int64) govalidity.ValidityResponseErrors {
	if files != nil {
		for _, file := range files {
			fileSize := file.Header.Size
			if fileSize > acceptedSizeInBytes {
				errs := govalidity.ValidityResponseErrors{}
				errs[label] = []string{
					fmt.Sprintf("Each file size should be less than %s", ConvertNumberToHumanReadableSize(acceptedSizeInBytes)),
				}
				return errs
			}
		}
	}
	return nil
}

// ValidateUploadFilesMimeType validates the mime type of the uploaded files
func ValidateUploadFilesMimeType(label string, files []*govalidityt.File, acceptedMimeTypes []string) govalidity.ValidityResponseErrors {
	if files != nil {
		for _, file := range files {
			fileMimeType := file.Header.Header.Get("Content-Type")
			if !Contains(acceptedMimeTypes, fileMimeType) {
				errs := govalidity.ValidityResponseErrors{}
				errs[label] = []string{
					fmt.Sprintf("Each file type should be one of %v", acceptedMimeTypes),
				}
				return errs
			}
		}
	}
	return nil
}
