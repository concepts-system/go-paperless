package domain

// DocumentMessageReceiver defines the callback for receiving a document.
type DocumentMessageReceiver = func(message DocumentMessage) error

// DocumentPageMessageReceiver defines the callback for receiving a document page.
type DocumentPageMessageReceiver = func(message DocumentPageMessage) error

// DocumentMessage defines the structure of a document related message in the tube mail system.
type DocumentMessage struct {
	DocumentNumber DocumentNumber
}

// DocumentPageMessage defines the structure of a document page related message in the tube mail system.
type DocumentPageMessage struct {
	DocumentNumber DocumentNumber
	PageNumber     PageNumber
}

// MailBox defines the type for a mail box.
type MailBox string

// TubeMail defines an interface for a tube mail system for moving documents
// between components.
type TubeMail interface {
	// RegisterDocumentMessageReceiver registeres a new document receiver for a given mail box.
	RegisterDocumentMessageReceiver(mailBox MailBox, receiver DocumentMessageReceiver) error

	// RegisterDocumentPageMessageReceiver registeres a new document receiver for a given mail box.
	RegisterDocumentPageMessageReceiver(mailBox MailBox, receiver DocumentPageMessageReceiver) error

	// SendDocumentMessage sends a document to a target mailbox.
	SendDocumentMessage(target MailBox, message DocumentMessage) error

	// SendDocumentPageMessage sends a document page to a target mailbox.
	SendDocumentPageMessage(target MailBox, message DocumentPageMessage) error
}
