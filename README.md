# Preguntas

*Utilizando buenas practicas de separación de capas, refactorizar handlers/handlers.go siguiendo una lógica de negocio/presentación
*Como utilizar correctamente el context en la solicitud
*Crear los test teniendo en cuenta las handlers_test.go
*Creame un readme a modo documentacion.

# Bookshop API

## Overview
The Bookshop API is a simple RESTful API built in Go using the [Gin framework](https://github.com/gin-gonic/gin). It retrieves book data from an external mock API (`https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books`) and provides metrics about the books, including:
- Average units sold across all books.
- Name of the cheapest book.
- Number of books written by a specified author.

The project follows a modular structure, separating concerns into models, providers, services, and handlers, with unit tests for key components.

## Project Structure
```
bookshop/
├── go.mod
├── go.sum
├── main.go
├── main_test.go
├── handlers/
│   └── handler.go
│   └── handler_test.go
├── models/
│   └── models.go
├── repositories/
│   ├── books.go
│   └── books_test.go
├── services/
│   ├── metrics.go
│   └── metrics_test.go
```

- `main.go`: Entry point, sets up the Gin server and routes.
- `handlers/`: Contains the request handler logic for processing API requests.
- `models/`: Defines the `Book` data structure.
- `repositories/`: Handles fetching book data from an external API.
- `services/`: Contains business logic for calculating book metrics.
- `*_test.go`: Unit tests for `providers` and `services` packages.

## Prerequisites
- Go 1.21 or higher (due to the use of `slices.MinFunc` in `services`).
- Git for cloning the repository.
- Internet access to fetch dependencies and the external mock API.

## Installation
1. **Clone the repository**:
   ```bash
   git clone https://github.com/freischarler/eb_challenge.git
   cd eb_challenge
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```
   This will download required packages, including:
   - `github.com/gin-gonic/gin` (for the web server)
   - `github.com/stretchr/testify` (for unit tests)

## Running the API
1. **Start the server**:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:3000`.

2. **Access the API**:
   - Endpoint: `GET /`
   - Optional query parameter: `author` (e.g., `?author=John%20Doe`)
   - Example request:
     ```bash
     curl "http://localhost:3000/?author=Ursula K. Le Guin
     ```
   - Example response:
     ```json
     {
        "mean_units_sold": 59333333,
        "cheapest_book": "A Wizard of Earthsea",
        "books_written_by_author": 1
     }
     ```

## API Details
- **Endpoint**: `GET /`
- **Query Parameter**:
  - `author` (string, optional): Filters the number of books by the specified author.
- **Response**:
  - `mean_units_sold` (uint): Average number of units sold across all books.
  - `cheapest_book` (string): Name of the book with the lowest price.
  - `books_written_by_author` (uint): Number of books by the specified author (0 if no author is provided or no books match).

## Error Handling

The API implements comprehensive error handling across all layers:

### Repository Layer Errors
- **`ErrServiceUnavailable`**: External service connection failed
- **Network timeouts**: 10-second timeout on HTTP requests
- **Invalid responses**: Non-200 HTTP status codes
- **JSON parsing errors**: Malformed external API responses

### Service Layer Errors  
- **`ErrExternalServiceFailure`**: Wraps repository errors for domain consistency
- **`ErrBookNotFound`**: No books available for processing

### Handler Layer Error Responses

| Scenario | HTTP Status | Response |
|----------|-------------|----------|
| Success | 200 OK | Metrics JSON |
| Invalid query parameters | 400 Bad Request | `{"error": "invalid query parameters"}` |
| External service failure | 502 Bad Gateway | `{"error": "error fetching books from external service"}` |
| Internal server error | 500 Internal Server Error | `{"error": "internal server error"}` |

**Error Response Format:**
```json
{
  "error": "error description"
}
```

## Running Tests
Unit tests are provided for the `providers` and `services` packages, using `testing` and `testify`.

1. **Run all tests**:
   ```bash
   go test ./... -v
   ```

2. **Run tests for a specific package**:
   - For `providers`:
     ```bash
     cd providers
     go test -v
     ```
   - For `services`:
     ```bash
     cd services
     go test -v
     ```

3. **Check test coverage**:
   ```bash
   go test ./... -cover
   ```
   For a detailed coverage report:
   ```bash
   go test ./... -coverprofile=cover.out && go tool cover -html=cover.out
   ```

## Example Data
The external API (`https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books`) may return data like:
```json
[
  {"id": 1, "name": "Book A", "author": "John Doe", "units_sold": 100, "price": 10},
  {"id": 2, "name": "Book B", "author": "John Doe", "units_sold": 200, "price": 15},
  {"id": 3, "name": "Book C", "author": "Jane Smith", "units_sold": 50, "price": 8}
]
```

A `GET /?author=John%20Doe` request would return:
```json
{
  "mean_units_sold": 116,
  "cheapest_book": "Book C",
  "books_written_by_author": 2
}
```

## Notes
- The external API URL is hardcoded in `main.go`. Consider using environment variables for flexibility.
- Logging is basic (uses `log.Printf`). Consider a structured logging library like `zerolog` for production.