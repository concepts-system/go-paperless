package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentContentKey(t *testing.T) {
	document := Document{
		Fingerprint: Fingerprint("fingerprint"),
		Type:        DocumentTypePDF,
	}

	assert.Equal(t, document.ContentKey(), ContentKey("fingerprint.pdf"))
}

func TestAllPagesAreInState(t *testing.T) {
	editedPage := DocumentPage{
		State: PageStateEdited,
	}

	analyzedPage := DocumentPage{
		State: PageStateAnalyzed,
	}

	partlyAnalyzedDocument := Document{
		Pages: []DocumentPage{editedPage, analyzedPage},
	}

	fullyAnalyzedDocument := Document{
		Pages: []DocumentPage{analyzedPage, analyzedPage},
	}

	emptyDocument := Document{}

	assert.Equal(t, partlyAnalyzedDocument.AreAllPagesInState(PageStateEdited), false)
	assert.Equal(t, partlyAnalyzedDocument.AreAllPagesInState(PageStateAnalyzed), false)
	assert.Equal(t, fullyAnalyzedDocument.AreAllPagesInState(PageStateEdited), false)
	assert.Equal(t, fullyAnalyzedDocument.AreAllPagesInState(PageStateAnalyzed), true)
	assert.Equal(t, emptyDocument.AreAllPagesInState(PageStateAnalyzed), true)
	assert.Equal(t, emptyDocument.AreAllPagesInState(PageStateEdited), true)
}
