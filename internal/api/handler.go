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
	// logger  *logrus.Logger
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

func (h handler) CreateAccount(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Invalid request body for creating account: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	if err := h.service.CreateAccount(c.Request.Context(), user.Email, user.Password); err != nil {
		log.Printf("Error creating account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "created")
}

func (h handler) GetOrderHistory(c *gin.Context) {

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}

	orders, err := h.service.GetOrderHistory(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error fetching order history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No Order found for this user"})
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

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, "email parameter is required")
		return
	}

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

	if err := h.service.PlaceOrder(c.Request.Context(), userID, books); err != nil {
		log.Printf("Error placing order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "order placed successfully")
}

func (h handler) GetUserIDByEmail(c *gin.Context) {

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, "email parameter is required")
		return
	}
	userID, err := h.service.GetUserIDByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}

	if userID == "" {
		c.JSON(http.StatusOK, gin.H{"message": "No user found for this email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
