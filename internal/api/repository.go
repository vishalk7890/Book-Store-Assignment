// package api

// import (
// 	"bookstore/db"
// 	"bookstore/internal/application"
// 	"context"
// )

// type Repository interface {
// 	GetAllBooks(ctx context.Context) ([]Book, error)
// 	PlaceOrder(ctx context.Context, email string, books []BookOrder) error
// 	CreateAccount(ctx context.Context, email string) error
// 	GetOrderHistory(ctx context.Context, email string) ([]Order, error)
// }

// type repository struct {
// 	app *application.Application
// 	db  *db.PostgresDB
// }

// func NewRepository(app *application.Application, dbRepo *db.PostgresDB) Repository {
// 	return &repository{
// 		app: app,
// 		db:  dbRepo,
// 	}
// }

// func (r *repository) CreateAccount(ctx context.Context, email string) error {

// 	return r.db.CreateAccount(ctx, email)
// }

// func (r *repository) GetAllBooks(ctx context.Context) ([]Book, error) {

// 	return r.db.GetAllBooks(ctx)

// }

// func (r *repository) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {

// 	// return r.db.PlaceOrder(ctx, email, books)
// 	return nil
// }

// func (r *repository) GetOrderHistory(ctx context.Context, email string) ([]Order, error) {

// 	return r.db.GetOrderHistory(ctx, email)
// }

package api

import (
	"bookstore/internal/application"
	"bookstore/internal/application/config"
	"context"
	"database/sql"
	"errors"
	"log"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
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
	CreateAccount(ctx context.Context, email string) error
	GetOrderHistory(ctx context.Context, email string) ([]Order, error)
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

func (r *repository) CreateAccount(ctx context.Context, email string) error {

	query := "INSERT INTO users (email) VALUES ($1)"
	_, err := r.db.db.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}
	return errors.New("not implemented")
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

func (r *repository) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {

	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order into the orders table
	orderID, err := r.insertOrder(ctx, tx, email)
	if err != nil {
		return err
	}

	// Insert order items into the order_items table
	for _, book := range books {
		if err := r.insertOrderItem(ctx, tx, orderID, book); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return errors.New("not implemented")
}

func (r *repository) GetOrderHistory(ctx context.Context, email string) ([]Order, error) {

	query := "SELECT id FROM orders WHERE user_id = $1"
	rows, err := r.db.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var orderID string
		if err := rows.Scan(&orderID); err != nil {
			log.Println("Error scanning order row:", err)
			continue
		}

		order, err := r.getOrder(ctx, orderID)
		if err != nil {
			log.Println("Error fetching order details:", err)
			continue
		}

		orders = append(orders, order)
	}

	return orders, errors.New("not implemented")
}

func (r repository) getOrder(ctx context.Context, orderID string) (Order, error) {
	var order Order
	query := "SELECT user_id, total_cost, status FROM orders WHERE id = $1"
	err := r.db.db.QueryRowContext(ctx, query, orderID).Scan(&order.UserID, &order.TotalCost)
	if err != nil {
		return Order{}, err
	}

	order.Items, err = r.getOrderItems(ctx, orderID)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (r repository) getOrderItems(ctx context.Context, orderID string) ([]OrderItem, error) {
	query := "SELECT book_id, price, quantity FROM order_items WHERE order_id = $1"
	rows, err := r.db.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.BookID, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r repository) insertOrder(ctx context.Context, tx *sql.Tx, email string) (string, error) {
	var orderID string
	query := "INSERT INTO orders (user_id) VALUES ($1) RETURNING id"
	err := tx.QueryRowContext(ctx, query, email).Scan(&orderID)
	if err != nil {
		return "", err
	}
	return orderID, nil

}
func (r repository) insertOrderItem(ctx context.Context, tx *sql.Tx, orderID string, book BookOrder) error {
	query := "INSERT INTO order_items (order_id, book_id, price, quantity) VALUES ($1, $2, $3, $4)"
	_, err := tx.ExecContext(ctx, query, orderID, book.BookID, book.Quantity)
	if err != nil {
		return err
	}
	return nil
}
