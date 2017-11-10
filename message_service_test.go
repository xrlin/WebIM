package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func setup() {
	redisConn := RedisPool.Get()
	defer func() {
		redisConn.Close()
	}()
	redisConn.Do("DEL", MessageList)
}

func TestPushMessageSuccess(t *testing.T) {
	msg := Message{RoomId: uint(rand.Uint32())}
	if err := PushMessage(msg); err != nil {
		t.Errorf("Push message to list failed with error: %s", err)
	}
}

func TestBRPopMessageSuccess(t *testing.T) {
	msg := Message{UUID: fmt.Sprintf("%v-%v", time.Now().Unix(), 1), RoomId: uint(rand.Uint32())}
	PushMessage(msg)
	replyMsg, err := BRPopMessage()
	if err != nil {
		t.Errorf("Pop message to list failed with error: %s", err)
	}
	if replyMsg.UUID != msg.UUID {
		t.Errorf("Poped message incorrect. Expected %v but %v instead!", msg, replyMsg)
	}
}

func TestMonitorAndDeliverMessages(t *testing.T) {
	hub := NewHub()
	msg := Message{UUID: fmt.Sprintf("%v-%v", time.Now().Unix(), 1), RoomId: uint(rand.Uint32())}
	PushMessage(msg)
	go MonitorAndDeliverMessages(hub)
	const waitDeliver = 3
	time.Sleep(waitDeliver)
	if len(hub.Messages) != 1 {
		t.Error("Fail to monitor and deliver message to hub!")
	}
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}
