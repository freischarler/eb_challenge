package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	// Arrange & Act
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Assert
	assert.NotNil(t, router)
}

func TestMain_GetMetrics_Integration(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := setupRouter()

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

	// This will test against the real external API
	assert.NotNil(t, result.MeanUnitsSold)
	assert.NotEmpty(t, result.CheapestBook)
}

func TestMain_GetMetrics_Integration_NoAuthor(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Act
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var result services.MetricsResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.NotNil(t, result.MeanUnitsSold)
	assert.NotEmpty(t, result.CheapestBook)
	assert.Equal(t, uint(0), result.BooksWrittenByAuthor)
}

func TestMain_RouteNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Act
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}
