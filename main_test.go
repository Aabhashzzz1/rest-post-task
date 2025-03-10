package main

import (
	"bytes"
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

func TestCreateUser(t *testing.T) {
	router := setupRouter()

	tests := []struct{ 
		name string 
		payload string 
		expected int 
	}{
		{
			name: "Invalid JSON",
			payload: `{name: "Inavlid JSON"}`,
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid Pan",
			payload: `[{
				"name": "aabhash",
				"pan": "12345ABCDE",
				"mobile": "9876543210",
				"email": "aabhashmalviya15@gmail.com"
			}]`,
			expected: http.StatusOK,
		},
		{
			name: "Invalid Mobile",
			payload: `[{
				"name": "aabhash",
				"pan": "ABCDE1234F",
				"mobile": "987654",
				"email": "aabhashmalviya15@gmail.com"
			}]`,
			expected: http.StatusOK,
		},
		{
			name: "Invalid Email",
			payload: `[{
				"name": "aabhash",
				"pan": "ABCDE1234F",
				"mobile": "987654",
				"email": "aabhash"
			}]`,
			expected: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func (t *testing.T) {
			req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expected {
				t.Errorf("[%s] Expected status %d, got %d\n", tc.name, tc.expected, w.Code)
			}
			t.Logf("[%s] Response: %s", tc.name, w.Body.String())
		})
	}
}