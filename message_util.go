package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
)

func MarshalMessage(msg Message) (msgText string, err error) {
	byteArray, err := json.Marshal(msg)
	msgText = string(byteArray)
	return msgText, err
}

func UnMarshalMessage(msgText string) (msg Message, err error) {
	err = json.Unmarshal([]byte(msgText), &msg)
	return msg, err
}

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
