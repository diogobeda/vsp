package messaging

import (
	"reflect"
	"testing"

	"github.com/nsqio/go-nsq"
)

type Message struct {
	Text string
}

func TestEncodeDecode(t *testing.T) {
	originalMessage := &Message{
		Text: "asdf",
	}
	encodedBody, encodeErr := EncodeMessage(originalMessage)
	if encodeErr != nil {
		t.Error("Encoder failed", encodeErr)
	}
	nsqMessage := &nsq.Message{
		Body: encodedBody,
	}

	decodedMessage := &Message{}
	decodeErr := DecodeMessage(nsqMessage, decodedMessage)
	if decodeErr != nil {
		t.Error("Decoder failed", decodeErr)
	}

	if !reflect.DeepEqual(decodedMessage, originalMessage) {
		t.Errorf("Expected original message to equal decoded message. %v != %v", originalMessage, decodedMessage)
	}
}
