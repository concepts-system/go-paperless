package infrastructure

import (
	"github.com/concepts-system/go-paperless/domain"
	log "github.com/sirupsen/logrus"
)

type receivers = map[domain.MailBox][]domain.TubeMailReceiver

type localAsyncTubeMailImpl struct {
	bufferSize int
	receivers  receivers
}

// NewTubeMailChannelImpl creates a new tube mail implementation using local channels.
func NewTubeMailChannelImpl() domain.TubeMail {
	return &localAsyncTubeMailImpl{
		bufferSize: 128,
		receivers:  make(receivers),
	}
}

func (t *localAsyncTubeMailImpl) RegisterReceiver(
	mailBox domain.MailBox,
	receiver domain.TubeMailReceiver,
) error {
	if _, ok := t.receivers[mailBox]; !ok {
		t.receivers[mailBox] = []domain.TubeMailReceiver{receiver}
	} else {
		t.receivers[mailBox] = append(t.receivers[mailBox], receiver)
	}

	return nil
}

func (t *localAsyncTubeMailImpl) SendMessage(
	target domain.MailBox,
	message interface{},
) error {
	receivers, ok := t.receivers[target]

	if !ok {
		log.Warnf("No receivers for mailbox '%s' registered.", target)
		return nil
	}

	for _, receiver := range receivers {
		go t.sendMessage(message, receiver)
	}

	return nil
}

func (t *localAsyncTubeMailImpl) sendMessage(
	message interface{},
	receiver domain.TubeMailReceiver,
) {
	err := receiver(message)

	if err != nil {
		log.Errorf("Error while handling document message: %s", err.Error())
	}
}
