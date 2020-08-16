package infrastructure

import (
	"github.com/concepts-system/go-paperless/domain"
	log "github.com/kpango/glg"
)

type (
	documentReceivers     = map[domain.MailBox][]domain.DocumentMessageReceiver
	documentPageReceivers = map[domain.MailBox][]domain.DocumentPageMessageReceiver
)

type tubeMailChannelImpl struct {
	bufferSize            int
	documentReceivers     documentReceivers
	documentPageReceivers documentPageReceivers
}

// NewTubeMailChannelImpl creates a new tube mail implementation using local channels.
func NewTubeMailChannelImpl() domain.TubeMail {
	return &tubeMailChannelImpl{
		bufferSize:            128,
		documentReceivers:     make(documentReceivers),
		documentPageReceivers: make(documentPageReceivers),
	}
}

func (t *tubeMailChannelImpl) RegisterDocumentMessageReceiver(
	mailBox domain.MailBox,
	receiver domain.DocumentMessageReceiver,
) error {
	if _, ok := t.documentReceivers[mailBox]; !ok {
		t.documentReceivers[mailBox] = []domain.DocumentMessageReceiver{receiver}
	} else {
		t.documentReceivers[mailBox] = append(t.documentReceivers[mailBox], receiver)
	}

	return nil
}

func (t *tubeMailChannelImpl) RegisterDocumentPageMessageReceiver(
	mailBox domain.MailBox,
	receiver domain.DocumentPageMessageReceiver,
) error {
	if _, ok := t.documentPageReceivers[mailBox]; !ok {
		t.documentPageReceivers[mailBox] = []domain.DocumentPageMessageReceiver{receiver}
	} else {
		t.documentPageReceivers[mailBox] = append(t.documentPageReceivers[mailBox], receiver)
	}

	return nil
}

func (t *tubeMailChannelImpl) SendDocumentMessage(
	target domain.MailBox,
	message domain.DocumentMessage,
) error {
	receivers, ok := t.documentReceivers[target]

	if !ok {
		log.Warnf("No receivers for mailbox '%s' registered.", target)
		return nil
	}

	for _, receiver := range receivers {
		go t.sendDocumentMessage(message, receiver)
	}

	return nil
}

func (t *tubeMailChannelImpl) SendDocumentPageMessage(
	target domain.MailBox,
	message domain.DocumentPageMessage,
) error {
	receivers, ok := t.documentPageReceivers[target]

	if !ok {
		log.Warnf("No receivers for mailbox '%s' registered.", target)
		return nil
	}

	for _, receiver := range receivers {
		go t.sendDocumentPageMessage(message, receiver)
	}

	return nil
}

func (t *tubeMailChannelImpl) sendDocumentMessage(
	message domain.DocumentMessage,
	receiver domain.DocumentMessageReceiver,
) {
	err := receiver(message)
	if err != nil {
		log.Errorf("Error while handling document message: %s", err.Error())
	}
}

func (t *tubeMailChannelImpl) sendDocumentPageMessage(
	message domain.DocumentPageMessage,
	receiver domain.DocumentPageMessageReceiver,
) {
	err := receiver(message)
	if err != nil {
		log.Errorf("Error while handling document page message: %s", err.Error())
	}
}
