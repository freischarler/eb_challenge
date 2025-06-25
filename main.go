package main

import (
	"fmt"

	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/repositories"
	"educabot.com/bookshop/services"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Books repository
	booksRepo := repositories.NewExternalBooksRepository("https://6781684b85151f714b0aa5db.mockapi.io/api/v1/books")

	// Servicio con lógica
	service := services.NewMetricsService(booksRepo)

	// Handler con dependencias
	handler := handlers.NewHandler(service)

	// Rutas
	router.GET("/", handler.GetMetrics)

	return router
}

func main() {
	router := setupRouter()

	// Iniciar servidor
	fmt.Println("🚀 Starting server on :3000")
	router.Run(":3000")
}
