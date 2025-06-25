package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories/mockImpls"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetMetrics_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockRepo := mockImpls.NewMockBooksRepositories()
	metricsService := services.NewMetricsService(mockRepo)
	handler := NewHandler(metricsService)

	router := gin.New()
	router.GET("/", handler.GetMetrics)

	// Act
	author := url.QueryEscape("Robert C. Martin")
	req := httptest.NewRequest(http.MethodGet, "/?author="+author, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var result services.MetricsResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Expected values based on mock data:
	// Books: Go Programming (5000), Clean Code (15000), Pragmatic Programmer (13000)
	// Mean: (5000 + 15000 + 13000) / 3 = 11000
	// Cheapest: Go Programming (40)
	// Books by Robert C. Martin: 1 (Clean Code)
	assert.Equal(t, uint(11000), result.MeanUnitsSold)
	assert.Equal(t, "The Go Programming Language", result.CheapestBook)
	assert.Equal(t, uint(1), result.BooksWrittenByAuthor)
}

func TestHandler_GetMetrics_Success_AuthorWithNoBooks(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockRepo := mockImpls.NewMockBooksRepositories()
	metricsService := services.NewMetricsService(mockRepo)
	handler := NewHandler(metricsService)

	router := gin.New()
	router.GET("/", handler.GetMetrics)

	// Act
	author := url.QueryEscape("Unknown Author")
	req := httptest.NewRequest(http.MethodGet, "/?author="+author, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var result services.MetricsResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, uint(11000), result.MeanUnitsSold)
	assert.Equal(t, "The Go Programming Language", result.CheapestBook)
	assert.Equal(t, uint(0), result.BooksWrittenByAuthor) // No books by this author
}

func TestHandler_GetMetrics_Success_NoAuthorParameter(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockRepo := mockImpls.NewMockBooksRepositories()
	metricsService := services.NewMetricsService(mockRepo)
	handler := NewHandler(metricsService)

	router := gin.New()
	router.GET("/", handler.GetMetrics)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var result services.MetricsResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, uint(11000), result.MeanUnitsSold)
	assert.Equal(t, "The Go Programming Language", result.CheapestBook)
	assert.Equal(t, uint(0), result.BooksWrittenByAuthor) // Empty author means no matches
}

func TestHandler_GetMetrics_ServiceFailure(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	// Create a mock that returns an error
	mockRepo := &mockErrorRepository{}
	metricsService := services.NewMetricsService(mockRepo)
	handler := NewHandler(metricsService)

	router := gin.New()
	router.GET("/", handler.GetMetrics)

	// Act
	req := httptest.NewRequest(http.MethodGet, "/?author=Test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error fetching books from external service", response["error"])
}

func TestNewHandler(t *testing.T) {
	// Arrange
	mockRepo := mockImpls.NewMockBooksRepositories()
	metricsService := services.NewMetricsService(mockRepo)

	// Act
	handler := NewHandler(metricsService)

	// Assert
	assert.NotNil(t, handler)
}

// Mock repository that returns an error for testing failure scenarios
type mockErrorRepository struct{}

func (m *mockErrorRepository) GetBooksProvider(ctx context.Context) ([]models.Book, error) {
	return nil, errors.New("repository error")
}
