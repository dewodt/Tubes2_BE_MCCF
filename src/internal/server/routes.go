package server

import (
	"net/http"
	"tubes2-be-mccf/internal/cache"
	"tubes2-be-mccf/internal/controllers"
	"tubes2-be-mccf/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() http.Handler {
	// Initialize gin router
	r := gin.Default()

	// CORS settings middleware
	r.Use(middleware.CORS())

	// Initialize cache json file
	cache.InitCache()

	// Endpoint for calculating wikirace shortest path
	r.POST("/play", controllers.PlayHandler)

	return r
}
