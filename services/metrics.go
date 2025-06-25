package services

import (
	"context"
	"errors"
	"slices"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories"
)

var ErrExternalServiceFailure = errors.New("error fetching books from external service")
var ErrBookNotFound = errors.New("book not found")

type MetricsResult struct {
	MeanUnitsSold        uint   `json:"mean_units_sold"`
	CheapestBook         string `json:"cheapest_book"`
	BooksWrittenByAuthor uint   `json:"books_written_by_author"`
}

type MetricsService struct {
	booksRepositories repositories.BooksRepository
}

func NewMetricsService(repository repositories.BooksRepository) *MetricsService {
	return &MetricsService{booksRepositories: repository}
}

func (s *MetricsService) ComputeMetrics(ctx context.Context, author string) (*MetricsResult, error) {
	books, err := s.booksRepositories.GetBooks(ctx)
	if err != nil {
		return nil, ErrExternalServiceFailure
	}

	result := &MetricsResult{
		MeanUnitsSold:        s.meanUnitsSold(books),
		CheapestBook:         s.cheapestBook(books).Name,
		BooksWrittenByAuthor: s.booksWrittenByAuthor(books, author),
	}
	return result, nil
}

func (s *MetricsService) meanUnitsSold(books []models.Book) uint {
	if len(books) == 0 {
		return 0
	}
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func (s *MetricsService) cheapestBook(books []models.Book) models.Book {
	if len(books) == 0 {
		return models.Book{}
	}
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func (s *MetricsService) booksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
