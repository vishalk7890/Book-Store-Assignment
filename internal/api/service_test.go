package api_test

import (
	"bookstore/internal/api"
	"bookstore/internal/api/mocks"
	"bookstore/internal/application"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Service_CreateAccount(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()
	tests := []struct {
		name        string
		email       string
		password    string
		repoErr     error
		expectedErr error
	}{
		{
			name:        "Successful account creation",
			email:       "test@example.com",
			password:    "password123",
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "Repository error",
			email:       "test@example.com",
			password:    "password123",
			repoErr:     errors.New("failed to create account, repository error"),
			expectedErr: errors.New("failed to create account, repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("CreateAccount", c, tt.email, tt.password).Return(tt.repoErr).Once()
			s := api.NewService(app, mockRepo)
			err := s.CreateAccount(c, tt.email, tt.password)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_Service_GetAllBooks(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()

	mockBooks := []api.Book{
		{ID: "1", Title: "Book 1", Author: "Author 1"},
		{ID: "2", Title: "Book 2", Author: "Author 2"},
	}

	tests := []struct {
		name          string
		repoBooks     []api.Book
		repoErr       error
		expectedBooks []api.Book
		expectedErr   error
	}{
		{
			name:          "Successful retrieval of books",
			repoBooks:     mockBooks,
			repoErr:       nil,
			expectedBooks: mockBooks,
			expectedErr:   nil,
		},
		{
			name:          "Repository error",
			repoBooks:     nil,
			repoErr:       errors.New("repository error"),
			expectedBooks: nil,
			expectedErr:   errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetAllBooks", c, mock.Anything).Return(tt.repoBooks, tt.repoErr).Once()
			s := api.NewService(app, mockRepo)

			books, err := s.GetAllBooks(context.Background())

			assert.Equal(t, tt.expectedBooks, books)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_Service_GetOrderHistory(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()

	mockUserID := "user1"
	mockOrders := []api.Order{
		{ID: "1", UserID: mockUserID, Items: []api.BookOrder{{BookID: "1", Quantity: 2, Title: "Book 1"}}},
		{ID: "2", UserID: mockUserID, Items: []api.BookOrder{{BookID: "2", Quantity: 1, Title: "Book 2"}}},
	}

	tests := []struct {
		name           string
		repoOrders     []api.Order
		repoUserID     string
		repoErr        error
		expectedOrders []api.Order
		expectedErr    error
	}{
		{
			name:           "Successful retrieval of order history",
			repoOrders:     mockOrders,
			repoUserID:     mockUserID,
			repoErr:        nil,
			expectedOrders: mockOrders,
			expectedErr:    nil,
		},
		{
			name:           "Repository error",
			repoOrders:     nil,
			repoUserID:     mockUserID,
			repoErr:        errors.New("repository error"),
			expectedOrders: nil,
			expectedErr:    errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetOrderHistory", c, tt.repoUserID).Return(tt.repoOrders, tt.repoErr).Once()
			svc := api.NewService(app, mockRepo)
			orders, err := svc.GetOrderHistory(context.Background(), tt.repoUserID)
			assert.Equal(t, tt.expectedOrders, orders)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_Service_PlaceOrder(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()

	mockEmail := "test@example.com"
	mockBooks := []api.BookOrder{
		{BookID: "1", Quantity: 2},
		{BookID: "2", Quantity: 1},
	}

	tests := []struct {
		name        string
		email       string
		books       []api.BookOrder
		repoErr     error
		expectedErr error
	}{
		{
			name:        "Successful order placement",
			email:       mockEmail,
			books:       mockBooks,
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:        "Repository error",
			email:       mockEmail,
			books:       mockBooks,
			repoErr:     errors.New("repository error"),
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("PlaceOrder", c, tt.email, tt.books).Return(tt.repoErr).Once()
			svc := api.NewService(app, mockRepo)
			err := svc.PlaceOrder(context.Background(), tt.email, tt.books)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_Service_GetBookByID(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()

	mockBookID := "1"
	mockBook := api.Book{
		ID:     mockBookID,
		Title:  "Test Book",
		Author: "Test Author",
	}

	tests := []struct {
		name         string
		bookID       string
		repoBook     api.Book
		repoErr      error
		expectedBook api.Book
		expectedErr  error
	}{
		{
			name:         "Successful retrieval of book",
			bookID:       mockBookID,
			repoBook:     mockBook,
			repoErr:      nil,
			expectedBook: mockBook,
			expectedErr:  nil,
		},
		{
			name:         "Repository error",
			bookID:       mockBookID,
			repoBook:     api.Book{},
			repoErr:      errors.New("repository error"),
			expectedBook: api.Book{},
			expectedErr:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetBookByID", c, tt.bookID).Return(tt.repoBook, tt.repoErr).Once()
			svc := api.NewService(app, mockRepo)
			book, err := svc.GetBookByID(context.Background(), tt.bookID)
			assert.Equal(t, tt.expectedBook, book)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func Test_Service_GetUserIDByEmail(t *testing.T) {
	app := application.NewAppMock()
	c := context.Background()

	mockEmail := "test@example.com"
	mockUserID := "1234"

	tests := []struct {
		name           string
		email          string
		repoUserID     string
		repoErr        error
		expectedUserID string
		expectedErr    error
	}{
		{
			name:           "Successful retrieval of user ID",
			email:          mockEmail,
			repoUserID:     mockUserID,
			repoErr:        nil,
			expectedUserID: mockUserID,
			expectedErr:    nil,
		},
		{
			name:           "Repository error",
			email:          mockEmail,
			repoUserID:     "",
			repoErr:        errors.New("repository error"),
			expectedUserID: "",
			expectedErr:    errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetUserIDByEmail", c, tt.email).Return(tt.repoUserID, tt.repoErr).Once()
			svc := api.NewService(app, mockRepo)
			userID, err := svc.GetUserIDByEmail(context.Background(), tt.email)
			assert.Equal(t, tt.expectedUserID, userID)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
