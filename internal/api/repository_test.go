package api_test

// import (
// 	"bookstore/internal/api"
// 	"bookstore/internal/application"
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockDB implements the PostgresDB interface for testing purposes.
// type MockPostgresDB struct {
// 	mock.Mock
// }

// func (m *MockPostgresDB) GetAllBooks(ctx context.Context) ([]api.Book, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]api.Book), args.Error(1)
// }

// func (m *MockPostgresDB) PlaceOrder(ctx context.Context, email string, books []api.BookOrder) error {
// 	args := m.Called(ctx, email, books)
// 	return args.Error(0)
// }

// // Implement other methods of the PostgresDB interface similarly

// func TestCreateAccount_Success(t *testing.T) {
// 	// Create a mock instance of PostgresDB
// 	mockDB := new(MockPostgresDB)
// 	// Assume the database operation succeeds without errors
// 	mockDB.On("CreateAccount", mock.Anything, mock.Anything, mock.Anything).Return(nil)

// 	// Create the repository with the mock database
// 	repo := api.NewRepository(&application.Application{}, mockDB)

// 	// Call the function under test
// 	err := repo.CreateAccount(context.Background(), "test@example.com", "password123")

// 	// Assert that no error occurred
// 	assert.NoError(t, err)

// 	// Assert that the CreateAccount method on the mock was called with the correct arguments
// 	mockDB.AssertCalled(t, "CreateAccount", mock.Anything, "test@example.com", "password123")
// }
