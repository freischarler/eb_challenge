package providers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"educabot.com/bookshop/models"
)

type BooksProvider interface {
	GetBooks(ctx context.Context) []models.Book
}

type ExternalBooksProvider struct {
	Endpoint string
}

func NewExternalBooksProvider(endpoint string) *ExternalBooksProvider {
	return &ExternalBooksProvider{Endpoint: endpoint}
}

func (p *ExternalBooksProvider) GetBooks(ctx context.Context) []models.Book {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.Endpoint, nil)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return []models.Book{}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error making HTTP request: %v", err)
		return []models.Book{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("non-200 response code: %d", resp.StatusCode)
		return []models.Book{}
	}

	var books []models.Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		log.Printf("error decoding response: %v", err)
		return []models.Book{}
	}

	return books
}
