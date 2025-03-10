package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(Logger())

	validate = validator.New()
	validate.RegisterValidation("pan", panValidator)
	validate.RegisterValidation("mobile", mobileValidator)

	r.POST("/users", CreateUsers)
	return r
}

func TestCreateUsers(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		payload    string
		statusCode int
		verify     func(t *testing.T, body string)
	}{
		{
			name:       "Invalid JSON",
			payload:    `{name: "Invalid JSON"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Empty Payload",
			payload:    `[]`,
			statusCode: http.StatusOK,
			verify: func(t *testing.T, body string) {
				var resp map[string]interface{}
				json.Unmarshal([]byte(body), &resp)
				if resp["success_count"].(float64) != 0 {
					t.Errorf("Expected success_count 0, got %v", resp["success_count"])
				}
			},
		},
		{
			name: "Valid User",
			payload: `[{
				"name": "Aabhash",
				"pan": "ABCDE1234F",
				"mobile": "9876543210",
				"email": "aabhash@example.com"
			}]`,
			statusCode: http.StatusOK,
			verify: func(t *testing.T, body string) {
				var resp map[string]interface{}
				json.Unmarshal([]byte(body), &resp)
				if resp["success_count"].(float64) != 1 {
					t.Errorf("Expected success_count 1, got %v", resp["success_count"])
				}
			},
		},
		{
			name: "Invalid PAN",
			payload: `[{
				"name": "Aabhash",
				"pan": "12345ABCDE",
				"mobile": "9876543210",
				"email": "aabhash@example.com"
			}]`,
			statusCode: http.StatusOK,
			verify: func(t *testing.T, body string) {
				var resp map[string]interface{}
				json.Unmarshal([]byte(body), &resp)
				if resp["failed_count"].(float64) != 1 {
					t.Errorf("Expected failed_count 1, got %v", resp["failed_count"])
				}
			},
		},
		{
			name: "Invalid Mobile",
			payload: `[{
				"name": "aabhash",
				"pan": "ABCDE1234F",
				"mobile": "987654",
				"email": "aabhash@example.com"
			}]`,
			statusCode: http.StatusOK,
		},
		{
			name: "Invalid Email",
			payload: `[{
				"name": "aabhash",
				"pan": "ABCDE1234F",
				"mobile": "9876543210",
				"email": "aabhash"
			}]`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.statusCode {
				t.Errorf("[%s] Expected status %d, got %d\n", tc.name, tc.statusCode, w.Code)
			}

			if tc.verify != nil {
				tc.verify(t, w.Body.String())
			}
		})
	}
}