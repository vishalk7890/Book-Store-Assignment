// package db

// import (
// 	"bookstore/internal/api"
// 	"bookstore/internal/application/config"
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"log"
// )

// type PostgresDB struct {
// 	db *sql.DB
// }

// func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
// 	db, err := sql.Open("postgres", cfg.DBConnectionString())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &PostgresDB{
// 		db: db,
// 	}, nil
// }

// func (pg *PostgresDB) GetAllBooks(ctx context.Context) ([]api.Book, error) {
// 	query := "SELECT id, title, author, description, price FROM books"
// 	rows, err := pg.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var books []api.Book
// 	for rows.Next() {
// 		var book api.Book
// 		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.Price)
// 		if err != nil {
// 			log.Println("Error scanning book row:", err)
// 			continue
// 		}
// 		books = append(books, book)
// 	}
// 	return books, nil
// }

// func (pg *PostgresDB) insertOrder(ctx context.Context, tx *sql.Tx, email string) (string, error) {
// 	var orderID string
// 	query := "INSERT INTO orders (user_id) VALUES ($1) RETURNING id"
// 	err := tx.QueryRowContext(ctx, query, email).Scan(&orderID)
// 	if err != nil {
// 		return "", err
// 	}
// 	return orderID, nil

// }
// func (pg *PostgresDB) insertOrderItem(ctx context.Context, tx *sql.Tx, orderID string, book api.BookOrder) error {
// 	query := "INSERT INTO order_items (order_id, book_id, price, quantity) VALUES ($1, $2, $3, $4)"
// 	_, err := tx.ExecContext(ctx, query, orderID, book.BookID, book.Quantity)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (pg *PostgresDB) PlaceOrder(ctx context.Context, email string, books []api.BookOrder) error {
// 	tx, err := pg.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Insert order into the orders table
// 	orderID, err := pg.insertOrder(ctx, tx, email)
// 	if err != nil {
// 		return err
// 	}

// 	// Insert order items into the order_items table
// 	for _, book := range books {
// 		if err := pg.insertOrderItem(ctx, tx, orderID, book); err != nil {
// 			return err
// 		}
// 	}

// 	// Commit the transaction
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}

// 	return errors.New("not implemented")
// }

// func (pg *PostgresDB) CreateAccount(ctx context.Context, email string) error {
// 	query := "INSERT INTO users (email) VALUES ($1)"
// 	_, err := pg.db.ExecContext(ctx, query, email)
// 	if err != nil {
// 		return err
// 	}
// 	return errors.New("not implemented")
// }

// func (pg *PostgresDB) GetOrderHistory(ctx context.Context, email string) ([]api.Order, error) {
// 	query := "SELECT id FROM orders WHERE user_id = $1"
// 	rows, err := pg.db.QueryContext(ctx, query, email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var orders []api.Order
// 	for rows.Next() {
// 		var orderID string
// 		if err := rows.Scan(&orderID); err != nil {
// 			log.Println("Error scanning order row:", err)
// 			continue
// 		}

// 		order, err := pg.getOrder(ctx, orderID)
// 		if err != nil {
// 			log.Println("Error fetching order details:", err)
// 			continue
// 		}

// 		orders = append(orders, order)
// 	}

// 	return orders, errors.New("not implemented")

// }

// func (pg *PostgresDB) getOrder(ctx context.Context, orderID string) (api.Order, error) {
// 	var order api.Order
// 	query := "SELECT user_id, total_cost, status FROM orders WHERE id = $1"
// 	err := pg.db.QueryRowContext(ctx, query, orderID).Scan(&order.UserID, &order.TotalCost)
// 	if err != nil {
// 		return api.Order{}, err
// 	}

// 	order.Items, err = pg.getOrderItems(ctx, orderID)
// 	if err != nil {
// 		return api.Order{}, err
// 	}

// 	return order, nil
// }

// func (pg *PostgresDB) getOrderItems(ctx context.Context, orderID string) ([]api.OrderItem, error) {
// 	query := "SELECT book_id, price, quantity FROM order_items WHERE order_id = $1"
// 	rows, err := pg.db.QueryContext(ctx, query, orderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var items []api.OrderItem
// 	for rows.Next() {
// 		var item api.OrderItem
// 		if err := rows.Scan(&item.BookID, &item.Price, &item.Quantity); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, item)
// 	}

// 	return items, nil
// }
