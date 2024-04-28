package api

type Order struct {
	ID     string      `json:"id"`
	UserID string      `json:"userId"`
	Items  []BookOrder `json:"items"`
}

type OrderItem struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Book struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type BookOrder struct {
	BookID   string `json:"bookId"`
	Quantity int    `json:"quantity"`
	Title    string `json:"titles"`
}

type BookOrdersReq struct {
	BookOrder
	Book
}
