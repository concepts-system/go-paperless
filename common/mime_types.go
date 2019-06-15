package common

import (
	"mime"
)

// GetExtensionByMimeType returns the primary file extension associated with the given mime tpye.
// The extensions always starts with a '.' like ".html".
func GetExtensionByMimeType(mimeType string) (string, error) {
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	return extensions[0], nil
}
