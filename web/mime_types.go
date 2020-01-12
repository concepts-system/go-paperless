package web

import (
	"mime"
	"strings"
)

// GetExtensionByMimeType returns the primary file extension associated with the given mime tpye.
func GetExtensionByMimeType(mimeType string) (string, error) {
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	return strings.ToLower(extensions[0][1:]), nil
}
