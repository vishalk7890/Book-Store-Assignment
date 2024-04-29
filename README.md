# Bookstore API

This is a RESTful API for a bookstore application built using Go and the Gin web framework. It provides functionality for managing books, creating user accounts, placing orders, and viewing order history.

## Features
- Create and manage books
- Create user accounts
- Place orders for books
- View order history

## Prerequisites
- Go 1.16 or higher
- PostgreSQL database

## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/bookstore-api.git
    ```

2. Navigate to the project directory:
    ```bash
    cd bookstore-api
    ```

3. Install dependencies:
    ```bash
    go mod tidy
    ```

4. Set up the PostgreSQL database and update the database configuration in `internal/application/config/config.go`.

5. Build the application:
    ```bash
    make build
    ```

6. Run the application:
    ```bash
    make run
    ```

    The API will be available at [http://localhost:8080](http://localhost:8080).

## API Endpoints
- `GET /books`: Get all books
- `POST /accounts`: Create a new user account
- `POST /orders`: Place a new order
- `GET /order/history`: Get order history for the authenticated user
- `GET /users/:email`: Get user ID by email query parameter
- `GET /book_detail`: Get Book Details by bookID query paramter

## Testing
To run the tests:
```bash
make test
```
## Other Targets
Include other targets like:
- Staticchecks
- gosec
- lint
- govulncheck
  