package api

import (
	"bookstore/internal/application"
	"bookstore/internal/application/config"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	//db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()))
	db, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		return nil, err
	}
	return &PostgresDB{
		db: db,
	}, nil

}

type Repository interface {
	GetAllBooks(ctx context.Context) ([]Book, error)
	PlaceOrder(ctx context.Context, email string, books []BookOrder) error
	CreateAccount(ctx context.Context, email, password string) error
	GetOrderHistory(ctx context.Context, email string) ([]Order, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	GetBookByID(ctx context.Context, bookID string) (Book, error)
}

type repository struct {
	app *application.Application
	db  PostgresDB
}

func NewRepository(app *application.Application, dbRepo PostgresDB) Repository {
	return &repository{
		app: app,
		db:  dbRepo,
	}
}

func (r *repository) CreateAccount(ctx context.Context, email, password string) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2)"
	_, err := r.db.db.Exec(query, email, password)
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}
	return nil
}

func (r *repository) GetAllBooks(ctx context.Context) ([]Book, error) {
	query := "SELECT id, title, author, description, price FROM books"
	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)
		if err != nil {
			log.Println("Error scanning book row:", err)
			continue
		}
		books = append(books, book)
	}
	return books, nil

}
func (r *repository) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	var userID string
	query := "SELECT id FROM users WHERE email = $1"
	err := r.db.db.QueryRowContext(ctx, query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("email not found:%v", err)
		}
		return "", fmt.Errorf("failed to get user ID: %v", err)
	}
	return userID, nil
}

//	func (r *repository) PlaceOrder(ctx context.Context, userID string, books []BookOrder) error {
//		for _, book := range books {
//			query := "INSERT INTO orders (user_id, book_id, quantity) VALUES ($1, $2, $3)"
//			_, err := r.db.db.ExecContext(ctx, query, userID, book.BookID, book.Quantity)
//			if err != nil {
//				return fmt.Errorf("failed to place order: %v", err)
//			}
//		}
//		return nil
//	}
func (r *repository) PlaceOrder(ctx context.Context, userID string, books []BookOrder) error {
	// Start a transaction
	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert order into the 'orders' table
	var orderID int
	query := "INSERT INTO orders (user_id) VALUES ($1) RETURNING id"
	err = tx.QueryRowContext(ctx, query, userID).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	// Insert order items into the 'order_items' table
	for _, book := range books {
		query = "INSERT INTO order_items (order_id, book_id, quantity) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, query, orderID, book.BookID, book.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// func (r *repository) GetOrderHistory(ctx context.Context, userID string) ([]Order, error) {
// 	query := "SELECT id, user_id FROM orders WHERE user_id = $1"
// 	rows, err := r.db.db.QueryContext(ctx, query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var orders []Order
// 	for rows.Next() {
// 		var order Order
// 		err := rows.Scan(&order.ID, &order.UserID)
// 		if err != nil {
// 			log.Println("Error scanning order row:", err)
// 			continue
// 		}
// 		orders = append(orders, order)
// 	}
// 	return orders, nil
// }

// func (r *repository) GetOrderHistory(ctx context.Context, userID string) ([]Order, error) {
// 	// Define the SQL query to fetch order details and associated items
// 	query := `
//         SELECT o.id, o.user_id, oi.book_id, oi.quantity
//         FROM orders o
//         JOIN order_items oi ON o.id = oi.order_id
//         WHERE o.user_id = $1
//     `

// 	rows, err := r.db.db.QueryContext(ctx, query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Map to store orders and their associated items
// 	orderMap := make(map[string]*Order)

// 	// Iterate over the query results
// 	for rows.Next() {
// 		var orderID, bookID string
// 		var quantity int
// 		// Scan the result into variables
// 		err := rows.Scan(&orderID, &userID, &bookID, &quantity)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// If the order doesn't exist in the map, create a new one
// 		if _, ok := orderMap[orderID]; !ok {
// 			orderMap[orderID] = &Order{
// 				ID:     orderID,
// 				UserID: userID,
// 				Items:  make([]BookOrder, 0),
// 			}
// 		}
// 		// Add the item to the order's items list
// 		orderMap[orderID].Items = append(orderMap[orderID].Items, BookOrder{
// 			BookID:   bookID,
// 			Quantity: quantity,
// 		})
// 	}

// 	// Convert the map to a slice of orders
// 	var orders []Order
// 	for _, order := range orderMap {
// 		orders = append(orders, *order)
// 	}

// 	return orders, nil
// }

func (r *repository) GetOrderHistory(ctx context.Context, userID string) ([]Order, error) {
	// Define the SQL query to fetch order details, item details, and book titles

	query := `
        SELECT o.id, o.user_id, oi.book_id, oi.quantity, b.title
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        JOIN books b ON oi.book_id = b.id
        WHERE o.user_id = $1
    `

	rows, err := r.db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store orders and their associated items
	orderMap := make(map[string]*Order)

	// Iterate over the query results
	for rows.Next() {
		var orderID, bookID, title string
		var quantity int
		// Scan the result into variables
		err := rows.Scan(&orderID, &userID, &bookID, &quantity, &title)
		if err != nil {
			return nil, err
		}
		// If the order doesn't exist in the map, create a new one
		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &Order{
				ID:     orderID,
				UserID: userID,
				Items:  make([]BookOrder, 0),
			}
		}
		// Add the order item to the order's items list
		orderMap[orderID].Items = append(orderMap[orderID].Items, BookOrder{
			BookID:   bookID,
			Quantity: quantity,
			Title:    title,
		})
	}

	// Convert the map to a slice of orders
	var orders []Order
	for _, order := range orderMap {
		orders = append(orders, *order) // Dereference the pointer before appending
	}

	return orders, nil
}

func (r *repository) GetBookByID(ctx context.Context, bookID string) (Book, error) {
	// Define the SQL query to fetch book details by ID
	query := "SELECT id, title, author, description, price FROM books WHERE id = $1"

	// Execute the SQL query and retrieve book details
	var book Book
	err := r.db.db.QueryRowContext(ctx, query, bookID).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)
	if err != nil {
		// Handle the case where the book is not found or an error occurs
		return Book{}, fmt.Errorf("failed to fetch book details: %v", err)
	}

	return book, nil
}
