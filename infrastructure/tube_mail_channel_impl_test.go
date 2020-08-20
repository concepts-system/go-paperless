package infrastructure

import (
	"testing"
	"time"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/stretchr/testify/assert"
)

const (
	receiveTimeout = 100 * time.Millisecond

	testMailBox  = domain.MailBox("mailbox")
	wrongMailBox = domain.MailBox("wrong")
)

type testMessage struct {
	field string
}

var (
	message = testMessage{
		field: "test",
	}
)

func TestSendAndReceiveMessage(t *testing.T) {
	tubeMail := NewTubeMailChannelImpl()
	correctMailbox := make(chan interface{})
	wrongMailbox := make(chan interface{})

	_ = tubeMail.RegisterReceiver(testMailBox, messageToChannelForwander(correctMailbox))
	_ = tubeMail.RegisterReceiver(wrongMailBox, messageToChannelForwander(wrongMailbox))

	_ = tubeMail.SendMessage(testMailBox, message)

	assertDocumentMessageReceived(t, correctMailbox, message)
	assertNoMessageReceived(t, wrongMailbox)
}

func messageToChannelForwander(channel chan interface{}) domain.TubeMailReceiver {
	return func(m interface{}) error {
		channel <- m
		return nil
	}
}

func assertDocumentMessageReceived(
	t *testing.T,
	channel chan interface{},
	expectedMessage interface{},
) {
	select {
	case message := <-channel:
		assert.Equal(t, expectedMessage, message)
	case <-time.After(receiveTimeout):
		t.Fatal("Did not receive the expected message within timeout")
	}
}

func assertNoMessageReceived(t *testing.T, channel chan interface{}) {
	select {
	case <-channel:
		t.Fatal("Did receive message unexpectedly within timeout")
	case <-time.After(receiveTimeout):
		break
	}
}
