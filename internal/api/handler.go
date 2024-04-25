package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetAllBooks(c *gin.Context)
	PlaceOrder(c *gin.Context)
	GetOrderHistory(c *gin.Context)
	CreateAccount(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h handler) GetAllBooks(c *gin.Context) {
	books, err := h.service.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch the books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h handler) PlaceOrder(c *gin.Context) {
	var orderRequet Order
	if err := c.ShouldBindJSON(&orderRequet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID := "placehoder_user_id"
	if err := h.service.PlaceOrder(c.Request.Context(), userID, []BookOrder{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	c.JSON(http.StatusCreated, "created")
}

func (h handler) CreateAccount(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.service.CreateAccount(c.Request.Context(), user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "created")
}

func (h handler) GetOrderHistory(c *gin.Context) {

	userID := "placeholder_user_id"

	orders, err := h.service.GetOrderHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
