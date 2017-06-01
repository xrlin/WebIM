package services

import (
	"testing"
	"time"
)

func TestTokenService_Generate(t *testing.T) {
	ts := TokenService{time.Hour * 3, "test"}
	_, err := ts.Generate(1, "test")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Test of  generating token passed.")
	}
}

func TestTokenService_Validate(t *testing.T) {
	ts := TokenService{time.Hour * 3, "test"}
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InRlc3QiLCJ1c2VySWQiOjF9.UwdKSum09_qnUlje5H-_QEahg3n6cNQwxieFzQHVZLc"
	_, err := ts.Validate(tokenString)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log("Test of parsing token passed.")
	}
}
