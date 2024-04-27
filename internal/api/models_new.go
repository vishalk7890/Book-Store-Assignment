package api

// Order represents an order entity.
type Order struct {
	ID     string      `json:"id"`
	UserID string      `json:"userId"`
	Items  []BookOrder `json:"items"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
}

// User represents a user entity.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Book represents a book entity.
type Book struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// BookOrder represents a book and the quantity ordered.
type BookOrder struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
	Title    string `json:"title"`
}

type BookOrdersReq struct {
	BookOrder
	Book
}