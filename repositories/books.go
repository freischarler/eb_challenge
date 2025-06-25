package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"educabot.com/bookshop/models"
)

var ErrServiceUnavailable = errors.New("external service failure")

type BooksRepository interface {
	GetBooksProvider(ctx context.Context) ([]models.Book, error)
}

type ExternalBooksRepository struct {
	Endpoint string
}

func NewExternalBooksRepository(endpoint string) *ExternalBooksRepository {
	return &ExternalBooksRepository{Endpoint: endpoint}
}

func (r *ExternalBooksRepository) GetBooksProvider(ctx context.Context) ([]models.Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external service returned status %d", resp.StatusCode)
	}

	var books []models.Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		return nil, err
	}

	return books, nil
}
