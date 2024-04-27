package api_test

import (
	"bookstore/internal/api"
	"bookstore/internal/api/mocks"
	"bookstore/internal/application"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateAccount(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()
	tests := []struct {
		name           string
		requestBodyStr string
		serviceError   error
		wantBody       string
		wantCode       int
	}{
		{
			name:           "happy case",
			requestBodyStr: `{"email": "test@example.com", "password": "password123"}`,
			serviceError:   nil,
			wantBody:       "\"created\"",
			wantCode:       http.StatusCreated,
		},
		{
			name:           "empty email case",
			requestBodyStr: `{"email": "", "password": "password123"}`,
			serviceError:   nil,
			wantBody:       "{\"error\":\"email is required\"}",
			wantCode:       http.StatusBadRequest,
		},
		{
			name:           "empty password case",
			requestBodyStr: `{"email": "test@example.com", "password": ""}`,
			serviceError:   nil,
			wantBody:       "{\"error\":\"password is required\"}",
			wantCode:       http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			requestBody := strings.NewReader(tt.requestBodyStr)
			mockService := new(mocks.Service)

			mockService.On("CreateAccount", c, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(tt.serviceError).Once()

			r.POST("/accounts", api.NewHandler(app, mockService).CreateAccount)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/accounts", requestBody)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func Test_GetAllBooks(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()
	tests := []struct {
		name         string
		serviceBooks []api.Book
		serviceError error
		wantBody     string
		wantCode     int
	}{
		{
			name:         "success case",
			serviceBooks: []api.Book{{ID: "1", Title: "Book 1", Author: "Author 1", Description: "test", Price: 1234}},
			serviceError: nil,
			wantBody:     `[{"id":"1","title":"Book 1","author":"Author 1","description":"test","price":1234}]`,
			wantCode:     http.StatusOK,
		},
		{
			name:         "error case",
			serviceBooks: nil,
			serviceError: errors.New("failed to fetch books"),
			wantBody:     `{"error":"failed to fetch the books"}`,
			wantCode:     http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := gin.Default()

			mockService := new(mocks.Service)

			mockService.On("GetAllBooks", c).Return(tt.serviceBooks, tt.serviceError).Once()
			r.GET("/books", api.NewHandler(app, mockService).GetAllBooks)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/books", nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())

		})
	}
}

func Test_GetOrderHistory(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()
	tests := []struct {
		name         string
		email        string
		userID       string
		serviceError error
		wantBody     string
		wantCode     int
	}{
		{
			name:         "success_case",
			email:        "test@example.com",
			userID:       "user123",
			serviceError: nil,
			wantBody:     `[{"id":"123","userId":"user123","items":[{"bookId":"1","quantity":2,"title":""}]}]`,
			wantCode:     http.StatusOK,
		},
		{
			name:         "email_not_provided",
			email:        "",
			userID:       "",
			serviceError: nil,
			wantBody:     `{"error":"email parameter is required"}`,
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "invalid_email",
			email:        "invalid_email",
			userID:       "",
			serviceError: errors.New("error getting user ID"),
			wantBody:     `{"error":"failed to get user ID"}`,
			wantCode:     http.StatusInternalServerError,
		},
		{
			name:         "error_fetching_order_history",
			email:        "test@example.com",
			userID:       "",
			serviceError: errors.New("failed to get user ID"),
			wantBody:     `{"error":"failed to get user ID"}`,
			wantCode:     http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			mockService := new(mocks.Service)
			mockService.On("GetUserIDByEmail", c, tt.email).Return(tt.userID, tt.serviceError).Once()
			mockService.On("GetOrderHistory", c, tt.userID).Return([]api.Order{
				{ID: "123", UserID: tt.userID, Items: []api.BookOrder{{BookID: "1", Quantity: 2, Title: ""}}},
			}, nil).Once()
			r.GET("/orders", api.NewHandler(app, mockService).GetOrderHistory)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/orders?email="+tt.email, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func Test_PlaceOrder(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()
	tests := []struct {
		name         string
		requestBody  interface{}
		email        string
		userID       string
		serviceError error
		wantBody     string
		wantCode     int
	}{
		{
			name: "success_case",
			requestBody: api.Order{
				Items: []api.BookOrder{
					{BookID: "1", Quantity: 2},
					{BookID: "2", Quantity: 1},
				},
			},
			email:        "test@example.com",
			userID:       "user123",
			serviceError: nil,
			wantBody:     `"order placed successfully"`,
			wantCode:     http.StatusCreated,
		},
		{
			name:         "invalid_request_body",
			requestBody:  "invalid", // invalid JSON request body
			email:        "test@example.com",
			userID:       "",
			serviceError: nil,
			wantBody:     `{"error":"invalid request body"}`,
			wantCode:     http.StatusBadRequest,
		},
		{
			name: "error_getting_user_id",
			requestBody: api.Order{
				Items: []api.BookOrder{
					{BookID: "1", Quantity: 2},
				},
			},
			email:        "test@example.com",
			userID:       "",
			serviceError: errors.New("failed to get user ID"),
			wantBody:     `{"error":"failed to get user ID"}`,
			wantCode:     http.StatusInternalServerError,
		},
		{
			name: "error_placing_order",
			requestBody: api.Order{
				Items: []api.BookOrder{
					{BookID: "1", Quantity: 2},
				},
			},
			email:        "test@example.com",
			userID:       "user123",
			serviceError: errors.New("failed to get user ID"),
			wantBody:     `{"error":"failed to get user ID"}`,
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			mockService := new(mocks.Service)
			mockService.On("GetUserIDByEmail", c, tt.email).Return(tt.userID, tt.serviceError).Once()
			mockService.On("PlaceOrder", c, tt.userID, mock.AnythingOfType("[]api.BookOrder")).Return(tt.serviceError).Once()
			r.POST("/orders", api.NewHandler(app, mockService).PlaceOrder)
			w := httptest.NewRecorder()

			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/orders?email="+tt.email, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			// Print out the actual response body for debugging
			fmt.Println("Actual response body:", w.Body.String())

			// Unmarshal the expected JSON response body into a map
			var expectedBody map[string]interface{}
			if err := json.Unmarshal([]byte(tt.wantBody), &expectedBody); err != nil {
				t.Errorf("error parsing expected body JSON: %v", err)
				return
			}

			// Compare the unmarshalled expected body to the actual response body
			var actualBody map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &actualBody); err != nil {
				t.Errorf("error parsing actual body JSON: %v", err)
				return
			}

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, expectedBody, actualBody)
		})
	}
}