package server

import (
	"net/http"
	"os"
	"tubes2-be-mccf/internal/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() http.Handler {
	r := gin.Default()

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FE_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * 3600,
	}))

	r.POST("/play", controllers.PlayHandler)

	return r
}
