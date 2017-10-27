package utils

import (
	"fmt"
	"github.com/xrlin/WebIM/server/models"
	"testing"
)

func TestMarshalMessageSuccess(t *testing.T) {
	msg := models.Message{RoomId: 1}
	if _, err := MarshalMessage(msg); err != nil {
		t.Errorf("Marshal %v failed with error: %s", msg, err)
	}
}

func TestUnMarshalMessage(t *testing.T) {
	roomId := 1
	msgText := fmt.Sprintf(`{"room_id": %d}`, roomId)
	msg, err := UnMarshalMessage(msgText)
	if err != nil {
		t.Error("UnMarshal failed!")
	}
	if msg.RoomId != uint(roomId) {
		t.Errorf("Struct unmarshaled incorrect!")
	}
}
