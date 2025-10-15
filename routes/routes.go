package routes

import (
	"example/web-service-gin/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Album routes
	albumRoutes := router.Group("/albums")
	{
		albumRoutes.GET("", controllers.GetAlbums)       // GET /albums
		albumRoutes.GET("/:id", controllers.GetAlbum)    // GET /albums/:id
		albumRoutes.POST("", controllers.CreateAlbum)    // POST /albums
		albumRoutes.PUT("/:id", controllers.UpdateAlbum) // PUT /albums/:id
		albumRoutes.DELETE("/:id", controllers.DeleteAlbum) // DELETE /albums/:id
	}
}
