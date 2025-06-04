package handlers

import (
	"net/http"

	"educabot.com/bookshop/providers"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewGetMetrics(booksProvider providers.BooksProvider) GetMetrics {
	return GetMetrics{
		booksProvider: booksProvider,
		service:       services.NewMetricsService(),
	}
}

type GetMetrics struct {
	booksProvider providers.BooksProvider
	service       *services.MetricsService
}

func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query GetMetricsRequest
		ctx.ShouldBindQuery(&query)

		books := h.booksProvider.GetBooks(ctx)

		meanUnitsSold := h.service.MeanUnitsSold(books)
		cheapestBook := h.service.CheapestBook(books).Name
		booksWrittenByAuthor := h.service.BooksWrittenByAuthor(books, query.Author)

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"mean_units_sold":         meanUnitsSold,
			"cheapest_book":           cheapestBook,
			"books_written_by_author": booksWrittenByAuthor,
		})
	}
}
