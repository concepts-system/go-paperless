package domain

// TubeMailReceiver defines the callback for receiving a document.
type TubeMailReceiver = func(message ...interface{}) error

// Mailbox defines the type for a mail box.
type Mailbox string

// TubeMail defines an interface for a tube mail system used for sending things
// around.
type TubeMail interface {
	// RegisterReceiver registers a new document receiver for a given mail box.
	RegisterReceiver(mailBox Mailbox, receiver TubeMailReceiver) error

	// SendMessage sends a message to a target mailbox.
	SendMessage(target Mailbox, message ...interface{}) error
}
