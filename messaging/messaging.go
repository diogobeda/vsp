package messaging

import (
	"bytes"
	"encoding/gob"

	"github.com/nsqio/go-nsq"
)

func EncodeMessage(e interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(e)

	return buf.Bytes(), err
}

func DecodeMessage(message *nsq.Message, e interface{}) error {
	buf := bytes.NewBuffer(message.Body)
	decodeErr := gob.NewDecoder(buf).Decode(e)

	return decodeErr
}
