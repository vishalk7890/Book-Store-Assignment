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
	// err := utitlity.ExecuteInsertQuery(ctx, r.db.db, query, email, password)
	// //return fmt.Errorf("failed to create account: %v", err)
	// if err != nil {
	// 	return fmt.Errorf("failed to create account: %v", err)
	// }
	// return nil
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
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows:%v", err)
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

func (r *repository) PlaceOrder(ctx context.Context, userID string, books []BookOrder) error {

	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	var orderID int
	query := "INSERT INTO orders (user_id) VALUES ($1) RETURNING id"
	err = tx.QueryRowContext(ctx, query, userID).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	for _, book := range books {
		query = "INSERT INTO order_items (order_id, book_id, quantity) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, query, orderID, book.BookID, book.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *repository) GetOrderHistory(ctx context.Context, userID string) ([]Order, error) {

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

	orderMap := make(map[string]*Order)

	for rows.Next() {
		var orderID, bookID, title string
		var quantity int
		err := rows.Scan(&orderID, &userID, &bookID, &quantity, &title)
		if err != nil {
			return nil, err
		}
		if _, ok := orderMap[orderID]; !ok {
			orderMap[orderID] = &Order{
				ID:     orderID,
				UserID: userID,
				Items:  make([]BookOrder, 0),
			}
		}
		orderMap[orderID].Items = append(orderMap[orderID].Items, BookOrder{
			BookID:   bookID,
			Quantity: quantity,
			Title:    title,
		})
	}

	var orders []Order
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *repository) GetBookByID(ctx context.Context, bookID string) (Book, error) {
	query := "SELECT id, title, author, description, price FROM books WHERE id = $1"

	var book Book
	err := r.db.db.QueryRowContext(ctx, query, bookID).Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)
	if err != nil {
		return Book{}, fmt.Errorf("failed to fetch book details: %v", err)
	}

	return book, nil
}
