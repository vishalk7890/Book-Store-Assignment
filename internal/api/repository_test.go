package api_test

// type mockDB struct {
// 	mock.Mock
// }

// func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
// 	argsMock := m.Called(ctx, query, args)
// 	return argsMock.Get(0).(sql.Result), argsMock.Error(1)
// }

// type mockPostgresDB struct {
// 	*mockDB
// }

// func (m *mockPostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
// 	// Implement BeginTx as needed for your test cases
// 	return &sql.Tx{}, nil
// }

// func TestCreateAccount(t *testing.T) {
// 	// Prepare mock
// 	dbMock := &mockPostgresDB{mockDB: &mockDB{}}
// 	// Create repository instance with mock
// 	repo := api.NewRepository(&application.Application{}, dbMock)

// 	// Set expectations
// 	dbMock.On("ExecContext", mock.Anything, "INSERT INTO users (email, password) VALUES ($1, $2)", "test@example.com", "password123").Return(nil, nil)

// 	// Call the function under test
// 	err := repo.CreateAccount(context.Background(), "test@example.com", "password123")
// 	assert.NoError(t, err)

// 	// Assert expectations were met
// 	dbMock.AssertExpectations(t)
// }
