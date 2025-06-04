// go
package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

func TestExternalBooksProvider_GetBooks_Success(t *testing.T) {
	expectedBooks := []models.Book{
		{ID: 1, Name: "Book One", Author: "Author A", Price: 10, UnitsSold: 100},
		{ID: 2, Name: "Book Two", Author: "Author B", Price: 12, UnitsSold: 200},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expectedBooks)
	}))
	defer server.Close()

	provider := NewExternalBooksProvider(server.URL)
	books := provider.GetBooks(context.Background())

	assert.Equal(t, expectedBooks, books)
}

func TestExternalBooksProvider_GetBooks_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := NewExternalBooksProvider(server.URL)
	books := provider.GetBooks(context.Background())

	assert.Empty(t, books)
}

func TestExternalBooksProvider_GetBooks_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not a json"))
	}))
	defer server.Close()

	provider := NewExternalBooksProvider(server.URL)
	books := provider.GetBooks(context.Background())

	assert.Empty(t, books)
}
