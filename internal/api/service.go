package api

import (
	"bookstore/internal/application"
	"context"
	"errors"
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
	books, err := s.repo.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s service) CreateAccount(ctx context.Context, email, password string) error {
	err := s.repo.CreateAccount(ctx, email, password)
	if err != nil {
		return errors.New("failed to create account")
	}
	return nil
}

func (s service) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {
	// userID, err := s.repo.GetUserIDByEmail(ctx, email)
	// if err != nil {
	// 	return fmt.Errorf("failed to get userID:%v", err)
	// }
	// if err := s.repo.PlaceOrder(ctx, userID, books); err != nil {
	// 	return fmt.Errorf("failed to place order: %v", err)
	// }

	return s.repo.PlaceOrder(ctx, email, books)
	//return nil
}

// func (s service) PlaceOrder(ctx context.Context, email string, books []BookOrder) error {
// 	userID, err := s.repo.GetUserIDByEmail(ctx, email)
// 	if err != nil {
// 		return fmt.Errorf("failed to get userID:%v", err)
// 	}
// 	if err := s.repo.PlaceOrder(ctx, userID, books); err != nil {
// 		return fmt.Errorf("failed to place order: %v", err)
// 	}

//		//return s.repo.PlaceOrder(ctx, email, books)
//		return nil
//	}
func (s service) GetOrderHistory(ctx context.Context, email string) ([]Order, error) {

	return s.repo.GetOrderHistory(ctx, email)
}

func (s service) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	return s.repo.GetUserIDByEmail(ctx, email)
}

func (s service) GetBookByID(ctx context.Context, bookID string) (Book, error) {
	return s.repo.GetBookByID(ctx, bookID)
}
