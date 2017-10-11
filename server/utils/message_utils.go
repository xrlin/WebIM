package utils

import (
	"encoding/json"
	"github.com/xrlin/WebIM/server/models"
)

func MarshalMessage(msg models.Message) (msgText string, err error) {
	byteArray, err := json.Marshal(msg)
	msgText = string(byteArray)
	return msgText, err
}

func UnMarshalMessage(msgText string) (msg models.Message, err error) {
	err = json.Unmarshal([]byte(msgText), &msg)
	return msg, err
}
