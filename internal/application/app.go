package application

import (
	"bookstore/internal/application/config"
	"errors"
	"os"
)

type Application struct {
	config *config.Config
}

func Load() (*Application, error) {
	app := Application{}
	config, err := config.Load()
	if err != nil {
		return &app, err

	}
	app.config = config
	return &app, nil
}

func NewAppMock() *Application {
	env := map[string]string{"PORT": "8080"}
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			_ = errors.New("load test case internal errror")
		}
	}
	app, err := Load()
	if err != nil {
		_ = errors.New("cannot load application")
		return nil
	}
	return app
}
