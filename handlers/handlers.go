package handlers

import (
	"net/http"

	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.MetricsService
}

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewHandler(service *services.MetricsService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMetrics(ctx *gin.Context) {
	var query GetMetricsRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	result, err := h.service.ComputeMetrics(ctx, query.Author)
	if err != nil {
		if err == services.ErrExternalServiceFailure {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
