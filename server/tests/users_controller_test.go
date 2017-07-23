package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xrlin/WebIM/server/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var server *gin.Engine = routes.RouterEngin()

func TestGetUserTokenSuccess(t *testing.T) {
	params := `{"user_name":"test", "password": "test"}`
	payload := strings.NewReader(params)

	request := httptest.NewRequest("POST", "/api/user/token", payload)
	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Log(responseRecorder.Code)
		t.Error("Test failed")
	}
}

func TestGetUserTokenFailed(t *testing.T) {
	params := `{"user_name":"test", "password": "&*(#(#)))%_$"}`
	payload := strings.NewReader(params)

	request := httptest.NewRequest("POST", "/api/user/token", payload)
	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code == http.StatusOK {
		t.Error("Test failed")
	}
}

func TestCreateUserHandlerSuccess(t *testing.T) {
	params := fmt.Sprintf(`{"user_name":"test%d", "password": "test"} `, time.Now().Unix())
	payload := strings.NewReader(params)

	request := httptest.NewRequest("POST", "/api/users", payload)
	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Test %s Failed", request.URL)
	}
}
