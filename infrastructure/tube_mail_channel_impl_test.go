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

var (
	testDocumentMessage = domain.DocumentMessage{
		DocumentNumber: domain.DocumentNumber(123),
	}

	testDocumentPageMessage = domain.DocumentPageMessage{
		DocumentNumber: domain.DocumentNumber(456),
		PageNumber:     domain.PageNumber(789),
	}
)

func TestSendAndReceiveDocumentMessage(t *testing.T) {
	tubeMail := NewTubeMailChannelImpl()
	correctMailbox := make(chan domain.DocumentMessage)
	wrongMailbox := make(chan domain.DocumentMessage)

	_ = tubeMail.RegisterDocumentMessageReceiver(testMailBox, forwardDocumentMessageToChannel(correctMailbox))
	_ = tubeMail.RegisterDocumentMessageReceiver(wrongMailBox, forwardDocumentMessageToChannel(wrongMailbox))

	_ = tubeMail.SendDocumentMessage(testMailBox, testDocumentMessage)

	assertReceiveDocumentMessage(t, correctMailbox, testDocumentMessage)
	assertNoDocumentMessageReceived(t, wrongMailbox)
}

func TestSendAndReceiveDocumentPageMessage(t *testing.T) {
	tubeMail := NewTubeMailChannelImpl()
	correctMailbox := make(chan domain.DocumentPageMessage)
	wrongMailbox := make(chan domain.DocumentPageMessage)

	_ = tubeMail.RegisterDocumentPageMessageReceiver(testMailBox, forwardDocumentPageMessageToChannel(correctMailbox))
	_ = tubeMail.RegisterDocumentPageMessageReceiver(wrongMailBox, forwardDocumentPageMessageToChannel(wrongMailbox))

	_ = tubeMail.SendDocumentPageMessage(testMailBox, testDocumentPageMessage)

	assertReceiveDocumentPageMessage(t, correctMailbox, testDocumentPageMessage)
	assertNoDocumentPageMessageReceived(t, wrongMailbox)
}

func forwardDocumentMessageToChannel(channel chan domain.DocumentMessage) domain.DocumentMessageReceiver {
	return func(m domain.DocumentMessage) error {
		channel <- m
		return nil
	}
}

func forwardDocumentPageMessageToChannel(channel chan domain.DocumentPageMessage) domain.DocumentPageMessageReceiver {
	return func(m domain.DocumentPageMessage) error {
		channel <- m
		return nil
	}
}

func assertReceiveDocumentMessage(t *testing.T, channel chan domain.DocumentMessage, expectedMessage domain.DocumentMessage) {
	select {
	case message := <-channel:
		assert.Equal(t, expectedMessage, message)
	case <-time.After(receiveTimeout):
		t.Fatal("Did not receive the expected message within timeout")
	}
}

func assertReceiveDocumentPageMessage(t *testing.T, channel chan domain.DocumentPageMessage, expectedMessage domain.DocumentPageMessage) {
	select {
	case message := <-channel:
		assert.Equal(t, expectedMessage, message)
	case <-time.After(receiveTimeout):
		t.Fatal("Did not receive the expected message within timeout")
	}
}

func assertNoDocumentMessageReceived(t *testing.T, channel chan domain.DocumentMessage) {
	select {
	case <-channel:
		t.Fatal("Did receive message unexpectedly within timeout")
	case <-time.After(receiveTimeout):
		break
	}
}

func assertNoDocumentPageMessageReceived(t *testing.T, channel chan domain.DocumentPageMessage) {
	select {
	case <-channel:
		t.Fatal("Did receive message unexpectedly within timeout")
	case <-time.After(receiveTimeout):
		break
	}
}
