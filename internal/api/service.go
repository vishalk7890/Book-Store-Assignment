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
	books, err := s.repo.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s service) CreateAccount(ctx context.Context, email, password string) error {

	return s.repo.CreateAccount(ctx, email, password)
}

func (s service) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {

	return s.repo.PlaceOrder(ctx, email, books)
}
func (s service) GetOrderHistory(ctx context.Context, email string) ([]Order, error) {

	return s.repo.GetOrderHistory(ctx, email)
}
