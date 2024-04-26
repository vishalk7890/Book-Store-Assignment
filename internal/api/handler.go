package api

import (
	"bookstore/internal/application"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAllBooks(c *gin.Context)
	PlaceOrder(c *gin.Context)
	GetOrderHistory(c *gin.Context)
	CreateAccount(c *gin.Context)
	GetUserIDByEmail(c *gin.Context)
}

type handler struct {
	app     *application.Application
	service Service
}

func NewHandler(app *application.Application, service Service) Handler {
	return &handler{
		app:     app,
		service: service,
	}
}

func (h handler) GetAllBooks(c *gin.Context) {
	books, err := h.service.GetAllBooks(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching books: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

//initial one

// func (h handler) PlaceOrder(c *gin.Context) {
// 	var orderRequet Order
// 	if err := c.ShouldBindJSON(&orderRequet); err != nil {
// 		log.Printf("Invalid request body for placing order: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
// 		return
// 	}

// 	userID := "placehoder_user_id"
// 	if err := h.service.PlaceOrder(c.Request.Context(), userID, []BookOrder{}); err != nil {
// 		log.Printf("Error placing order: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

// 		return
// 	}
// 	c.JSON(http.StatusCreated, "created")
// }

func (h handler) CreateAccount(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Invalid request body for creating account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.service.CreateAccount(c.Request.Context(), user.Email, user.Password); err != nil {
		log.Printf("Error creating account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "created")
}

// func (h handler) GetOrderHistory(c *gin.Context) {

// 	userID := "placeholder_user_id"

// 	orders, err := h.service.GetOrderHistory(c.Request.Context(), userID)
// 	if err != nil {
// 		log.Printf("Error fetching order history: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

//		c.JSON(http.StatusOK, orders)
//	}
// func (h handler) GetOrderHistory(c *gin.Context) {
// 	// Get the email from the query parameters
// 	email := c.Query("email")

// 	// Get the user ID from the email
// 	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
// 	if err != nil {
// 		log.Printf("Error getting user ID: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
// 		return
// 	}

// 	// Now you have the user ID, you can use it to get the order history
// 	orders, err := h.service.GetOrderHistory(c.Request.Context(), userID)
// 	if err != nil {
// 		log.Printf("Error fetching order history: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Retrieve book titles from the order history
// 	var bookTitles []string
// 	for _, order := range orders {
// 		for _, item := range order.Items {
// 			// Fetch book details based on BookID
// 			book, err := h.service.GetBookByID(c.Request.Context(), item.BookID)
// 			if err != nil {
// 				log.Printf("Error fetching book details: %v", err)
// 				continue
// 			}
// 			bookTitles = append(bookTitles, book.Title)
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"bookTitles": bookTitles})
// }

func (h handler) GetOrderHistory(c *gin.Context) {
	// Get the email from the query parameters
	email := c.Query("email")

	// Get the user ID from the email
	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}

	// Now you have the user ID, you can use it to get the order history
	orders, err := h.service.GetOrderHistory(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error fetching order history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h handler) PlaceOrder(c *gin.Context) {
	var orderRequest Order
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		log.Printf("Invalid request body for placing order: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get the email from the query parameters
	email := c.Query("email")

	// Get the user ID from the email
	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}

	var books []BookOrder
	for _, item := range orderRequest.Items {
		books = append(books, BookOrder{
			BookID:   item.BookID,
			Quantity: item.Quantity,
		})
	}
	// Now you have the user ID, you can use it to place the order
	if err := h.service.PlaceOrder(c.Request.Context(), userID, books); err != nil {
		log.Printf("Error placing order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "order placed successfully")
}

func (h handler) GetUserIDByEmail(c *gin.Context) {

	email := c.Query("email")
	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
