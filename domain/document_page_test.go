package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentPageContentKey(t *testing.T) {
	page := DocumentPage{
		Fingerprint: Fingerprint("fingerprint"),
		Type:        PageTypeTIFF,
	}

	assert.Equal(t, page.ContentKey(), ContentKey("fingerprint.tiff"))
}
