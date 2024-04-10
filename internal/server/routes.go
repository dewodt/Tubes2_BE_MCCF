package server

import (
	"net/http"
	"tubes2-be-mccf/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() http.Handler {
	r := gin.Default()

	r.POST("/play", controllers.PlayHandler)

	return r
}
