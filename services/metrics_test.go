package services

import (
	"math"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

func sampleBooks() []models.Book {
	return []models.Book{
		{ID: 1, Name: "Book A", Author: "Author X", Price: 10, UnitsSold: 100},
		{ID: 2, Name: "Book B", Author: "Author Y", Price: 5, UnitsSold: 200},
		{ID: 3, Name: "Book C", Author: "Author X", Price: 20, UnitsSold: 300},
	}
}

func TestMeanUnitsSold(t *testing.T) {
	service := NewMetricsService()

	t.Run("MultipleBooks", func(t *testing.T) {
		books := sampleBooks()
		assert.Equal(t, uint(200), service.MeanUnitsSold(books), "should calculate mean correctly")
	})

	t.Run("SingleBook", func(t *testing.T) {
		books := []models.Book{{ID: 1, UnitsSold: 500}}
		assert.Equal(t, uint(500), service.MeanUnitsSold(books), "should return units sold for single book")
	})

	t.Run("EmptyList", func(t *testing.T) {
		books := []models.Book{}
		assert.Equal(t, uint(0), service.MeanUnitsSold(books), "should return 0 for empty list")
	})

	t.Run("ZeroUnitsSold", func(t *testing.T) {
		books := []models.Book{{UnitsSold: 0}, {UnitsSold: 0}}
		assert.Equal(t, uint(0), service.MeanUnitsSold(books), "should handle zero units sold")
	})

	t.Run("MaxUint", func(t *testing.T) {
		books := []models.Book{{UnitsSold: math.MaxUint32}, {UnitsSold: math.MaxUint32}}
		assert.Equal(t, uint(math.MaxUint32), service.MeanUnitsSold(books), "should handle max uint values")
	})
}

func TestCheapestBook(t *testing.T) {
	service := NewMetricsService()

	t.Run("MultipleBooks", func(t *testing.T) {
		books := sampleBooks()
		cheapest := service.CheapestBook(books)
		assert.Equal(t, uint(5), cheapest.Price, "should find correct cheapest price")
		assert.Equal(t, "Book B", cheapest.Name, "should find correct cheapest book name")
	})

	t.Run("SingleBook", func(t *testing.T) {
		books := []models.Book{{Name: "Book A", Price: 10}}
		cheapest := service.CheapestBook(books)
		assert.Equal(t, uint(10), cheapest.Price, "should return single book price")
		assert.Equal(t, "Book A", cheapest.Name, "should return single book name")
	})

	t.Run("EmptyList", func(t *testing.T) {
		books := []models.Book{}
		cheapest := service.CheapestBook(books)
		assert.Equal(t, models.Book{}, cheapest, "should return empty book for empty list")
	})

	t.Run("SamePrice", func(t *testing.T) {
		books := []models.Book{
			{Name: "Book A", Price: 5},
			{Name: "Book B", Price: 5},
		}
		cheapest := service.CheapestBook(books)
		assert.Equal(t, uint(5), cheapest.Price, "should handle same price")
		assert.Contains(t, []string{"Book A", "Book B"}, cheapest.Name, "should return one of the cheapest books")
	})
}

func TestBooksWrittenByAuthor(t *testing.T) {
	service := NewMetricsService()

	t.Run("MultipleBooksByAuthor", func(t *testing.T) {
		books := sampleBooks()
		assert.Equal(t, uint(2), service.BooksWrittenByAuthor(books, "Author X"), "should count multiple books by Author X")
	})

	t.Run("SingleBookByAuthor", func(t *testing.T) {
		books := sampleBooks()
		assert.Equal(t, uint(1), service.BooksWrittenByAuthor(books, "Author Y"), "should count single book by Author Y")
	})

	t.Run("UnknownAuthor", func(t *testing.T) {
		books := sampleBooks()
		assert.Equal(t, uint(0), service.BooksWrittenByAuthor(books, "Unknown"), "should return 0 for unknown author")
	})

	t.Run("EmptyAuthor", func(t *testing.T) {
		books := sampleBooks()
		assert.Equal(t, uint(0), service.BooksWrittenByAuthor(books, ""), "should return 0 for empty author")
	})

	t.Run("EmptyList", func(t *testing.T) {
		books := []models.Book{}
		assert.Equal(t, uint(0), service.BooksWrittenByAuthor(books, "Author X"), "should return 0 for empty list")
	})

	t.Run("CaseSensitivity", func(t *testing.T) {
		books := []models.Book{{Author: "Author X"}, {Author: "author x"}}
		assert.Equal(t, uint(1), service.BooksWrittenByAuthor(books, "Author X"), "should be case sensitive")
	})
}
