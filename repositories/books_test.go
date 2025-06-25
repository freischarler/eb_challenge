package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

func TestExternalBooksRepository_GetBooks_Success(t *testing.T) {
	// Arrange
	mockBooks := []models.Book{
		{ID: 1, Name: "Test Book", Author: "Test Author", UnitsSold: 100, Price: 25},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockBooks)
	}))
	defer server.Close()

	repo := NewExternalBooksRepository(server.URL)
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, mockBooks, books)
}

func TestExternalBooksRepository_GetBooks_InvalidURL(t *testing.T) {
	// Arrange
	repo := NewExternalBooksRepository("://invalid-url")
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, books)
}

func TestExternalBooksRepository_GetBooks_NetworkError(t *testing.T) {
	// Arrange
	repo := NewExternalBooksRepository("http://localhost:99999")
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.Equal(t, ErrServiceUnavailable, err)
	assert.Nil(t, books)
}

func TestExternalBooksRepository_GetBooks_NonOKStatus(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := NewExternalBooksRepository(server.URL)
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "external service returned status 500")
	assert.Nil(t, books)
}

func TestExternalBooksRepository_GetBooks_InvalidJSON(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	repo := NewExternalBooksRepository(server.URL)
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, books)
}

func TestExternalBooksRepository_GetBooks_ContextCancellation(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate slow response
		json.NewEncoder(w).Encode([]models.Book{})
	}))
	defer server.Close()

	repo := NewExternalBooksRepository(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.Equal(t, ErrServiceUnavailable, err)
	assert.Nil(t, books)
}

func TestExternalBooksRepository_GetBooks_EmptyResponse(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Book{})
	}))
	defer server.Close()

	repo := NewExternalBooksRepository(server.URL)
	ctx := context.Background()

	// Act
	books, err := repo.GetBooksProvider(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, books)
}

func TestNewExternalBooksRepository(t *testing.T) {
	// Arrange
	endpoint := "http://example.com/books"

	// Act
	repo := NewExternalBooksRepository(endpoint)

	// Assert
	assert.NotNil(t, repo)
	assert.Equal(t, endpoint, repo.Endpoint)
}
