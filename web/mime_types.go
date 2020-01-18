package web

import (
	"mime"
	"strings"
)

// getExtensionByMimeType returns the primary file extension associated with the given mime tpye.
func getExtensionByMimeType(mimeType string) (string, error) {
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	return strings.ToLower(extensions[0][1:]), nil
}
