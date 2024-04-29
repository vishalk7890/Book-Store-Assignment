package api

import (
	"bookstore/internal/application"
	"context"
)

type Service interface {
	GetAllBooks(ctx context.Context) ([]Book, error)
	CreateAccount(ctx context.Context, email, password string) error
	PlaceOrder(ctx context.Context, email string, books []BookOrder) error
	GetOrderHistory(ctx context.Context, email string) ([]Order, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	GetBookByID(ctx context.Context, bookID string) (Book, error)
}

type service struct {
	app  *application.Application
	repo Repository
}

func NewService(app *application.Application, repo Repository) Service {
	return &service{
		app:  app,
		repo: repo,
	}
}

func (s service) GetAllBooks(ctx context.Context) ([]Book, error) {

	return s.repo.GetAllBooks(ctx)
}

func (s service) CreateAccount(ctx context.Context, email, password string) error {

	return s.repo.CreateAccount(ctx, email, password)
}

func (s service) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {

	return s.repo.PlaceOrder(ctx, email, books)
	//return nil
}

func (s service) GetOrderHistory(ctx context.Context, email string) ([]Order, error) {

	return s.repo.GetOrderHistory(ctx, email)
}

func (s service) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	return s.repo.GetUserIDByEmail(ctx, email)
}

func (s service) GetBookByID(ctx context.Context, bookID string) (Book, error) {
	return s.repo.GetBookByID(ctx, bookID)
}
