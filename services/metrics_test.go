package services

import (
	"context"
	"errors"
	"testing"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories/mockImpls"
	"github.com/stretchr/testify/assert"
)

func TestMetricsService_ComputeMetrics_Success(t *testing.T) {
	// Arrange
	mockRepo := mockImpls.NewMockBooksRepositories()
	service := NewMetricsService(mockRepo)
	ctx := context.Background()
	author := "Robert C. Martin"

	// Act
	result, err := service.ComputeMetrics(ctx, author)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(11000), result.MeanUnitsSold)                  // (5000 + 15000 + 13000) / 3
	assert.Equal(t, "The Go Programming Language", result.CheapestBook) // Price 40
	assert.Equal(t, uint(1), result.BooksWrittenByAuthor)               // Clean Code by Robert C. Martin
}

func TestMetricsService_ComputeMetrics_AuthorWithNoBooks(t *testing.T) {
	// Arrange
	mockRepo := mockImpls.NewMockBooksRepositories()
	service := NewMetricsService(mockRepo)
	ctx := context.Background()
	author := "Unknown Author"

	// Act
	result, err := service.ComputeMetrics(ctx, author)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(11000), result.MeanUnitsSold)
	assert.Equal(t, "The Go Programming Language", result.CheapestBook)
	assert.Equal(t, uint(0), result.BooksWrittenByAuthor)
}

func TestMetricsService_ComputeMetrics_AuthorWithMultipleBooks(t *testing.T) {
	// Arrange
	mockRepo := mockImpls.NewMockBooksRepositories()
	service := NewMetricsService(mockRepo)
	ctx := context.Background()
	author := "Alan Donovan"

	// Act
	result, err := service.ComputeMetrics(ctx, author)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.BooksWrittenByAuthor)
}

func TestMetricsService_ComputeMetrics_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := &MockBooksRepositoryWithError{}
	service := NewMetricsService(mockRepo)
	ctx := context.Background()
	author := "Any Author"

	// Act
	result, err := service.ComputeMetrics(ctx, author)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrExternalServiceFailure, err)
	assert.Nil(t, result)
}

func TestNewMetricsService(t *testing.T) {
	// Arrange
	mockRepo := mockImpls.NewMockBooksRepositories()

	// Act
	service := NewMetricsService(mockRepo)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.booksRepositories)
}

func TestMetricsService_meanUnitsSold(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{
		{UnitsSold: 1000},
		{UnitsSold: 2000},
		{UnitsSold: 3000},
	}

	// Act
	result := service.meanUnitsSold(books)

	// Assert
	assert.Equal(t, uint(2000), result)
}

func TestMetricsService_meanUnitsSold_EmptySlice(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{}

	// Act
	result := service.meanUnitsSold(books)

	// Assert
	assert.Equal(t, uint(0), result)
}

func TestMetricsService_cheapestBook(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{
		{Name: "Expensive Book", Price: 100},
		{Name: "Cheap Book", Price: 20},
		{Name: "Medium Book", Price: 50},
	}

	// Act
	result := service.cheapestBook(books)

	// Assert
	assert.Equal(t, "Cheap Book", result.Name)
	assert.Equal(t, uint(20), result.Price)
}

func TestMetricsService_cheapestBook_EmptySlice(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{}

	// Act
	result := service.cheapestBook(books)

	// Assert
	assert.Equal(t, models.Book{}, result)
}

func TestMetricsService_booksWrittenByAuthor(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{
		{Author: "John Doe"},
		{Author: "Jane Smith"},
		{Author: "John Doe"},
		{Author: "Bob Wilson"},
	}

	// Act
	result := service.booksWrittenByAuthor(books, "John Doe")

	// Assert
	assert.Equal(t, uint(2), result)
}

func TestMetricsService_booksWrittenByAuthor_NoMatches(t *testing.T) {
	// Arrange
	service := &MetricsService{}
	books := []models.Book{
		{Author: "John Doe"},
		{Author: "Jane Smith"},
	}

	// Act
	result := service.booksWrittenByAuthor(books, "Unknown Author")

	// Assert
	assert.Equal(t, uint(0), result)
}

// Mock repository that returns an error for testing error scenarios
type MockBooksRepositoryWithError struct{}

func (m *MockBooksRepositoryWithError) GetBooks(ctx context.Context) ([]models.Book, error) {
	return nil, errors.New("repository error")
}
