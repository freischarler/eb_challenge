package services

import (
	"slices"

	"educabot.com/bookshop/models"
)

type MetricsService struct{}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

func (s *MetricsService) MeanUnitsSold(books []models.Book) uint {
	if len(books) == 0 {
		return 0
	}
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func (s *MetricsService) CheapestBook(books []models.Book) models.Book {
	if len(books) == 0 {
		return models.Book{}
	}
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func (s *MetricsService) BooksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
