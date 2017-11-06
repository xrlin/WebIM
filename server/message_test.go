package main

import (
	"fmt"
	"testing"
)

func TestMessage_RoomName(t *testing.T) {
	msg := Message{RoomId: 3}
	if msg.RoomName() != "room_3" {
		t.Fail()
	}
}

func TestMessage_UserRoomName(t *testing.T) {
	msg := Message{UserId: 3}
	fmt.Println(msg.User)
	if msg.UserRoomName() != "room_user_3" {
		t.Fail()
	}
}

func getUUID() string {
	uuid, _ := GenerateUUID()
	return uuid
}

func TestCreateMessage(t *testing.T) {
	type testCase struct {
		message   Message
		isSuccess bool
	}
	testCases := []testCase{
		{Message{UUID: "6a8f5e8f-eac6-4962-aebf-cf04751162e2", UserId: 1, MsgType: SingleMessage}, true},
		// UUID should be unique
		{Message{UUID: "6a8f5e8f-eac6-4962-aebf-cf04751162e2", UserId: 1, MsgType: SingleMessage}, false},
		{Message{UUID: "", UserId: 1, MsgType: SingleMessage}, false},
		{Message{UUID: getUUID(), UserId: 0, MsgType: SingleMessage}, false},
		{Message{UUID: getUUID(), UserId: 1, MsgType: 0}, false},
	}
	messages := make([]Message, len(testCases))
	for caseIndex, testCase := range testCases {
		message := testCase.message
		err := db.Create(&message).Error
		if err == nil && !testCase.isSuccess {
			t.Errorf("Message with %#v should failed \n test case %d", message, caseIndex)
		}
		if err != nil && testCase.isSuccess {
			t.Errorf("Message with %#v should success \n with err %s \n test case %d", message, err.Error(), caseIndex)
		}
		messages = append(messages, message)
	}
	defer func() {
		ids := []uint{}
		for _, msg := range messages {
			ids = append(ids, msg.ID)
		}
		db.Unscoped().Where("id IN (?)", ids).Delete(Message{})
	}()
}

func TestMessage_GetDetails(t *testing.T) {
	// init test data
	var user1, user2 User
	user1 = User{Avatar: "test", Name: randomUserName(), PasswordHash: getUUID()}
	user2 = User{Avatar: "test2", Name: randomUserName(), PasswordHash: getUUID()}
	db.Create(&user1)
	db.Create(&user2)
	defer func() {
		db.Delete(&user1)
		db.Delete(&user2)
	}()

	message := Message{UUID: getUUID(), UserId: user1.ID, FromUser: user2.ID}
	messageDetails := message.GetDetails()
	expected := MessageDetail{message, user2, user1}
	if (messageDetails.SourceUser.ID != expected.FromUser) && (messageDetails.TargetUser.ID != expected.TargetUser.ID) {
		t.Errorf("GetDetails result expected %#v but get %#v", expected, messageDetails)
	}
}
