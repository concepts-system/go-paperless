package domain

// TubeMailReceiver defines the callback for receiving a document.
type TubeMailReceiver = func(message interface{}) error

// MailBox defines the type for a mail box.
type MailBox string

// TubeMail defines an interface for a tube mail system used for sending things
// around.
type TubeMail interface {
	// RegisterReceiver registeres a new document receiver for a given mail box.
	RegisterReceiver(mailBox MailBox, receiver TubeMailReceiver) error

	// SendMessage sends a message to a target mailbox.
	SendMessage(target MailBox, message interface{}) error
}
