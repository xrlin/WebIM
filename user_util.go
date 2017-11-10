package main

import (
	"fmt"
	"time"
)

func randomUserName() string {
	return fmt.Sprintf("User%d", time.Now().Unix())
}
