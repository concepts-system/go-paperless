package domain

import "io"

// DocumentGenerator defines a signature for a component being able to generate
// actual, human-readable documents.
type DocumentGenerator interface {
	// Generate generates the given document and returns a reader
	// for the generated content.
	Generate(document *Document) (io.Reader, error)
}
