package routes

import (
	"web-service-gin/backend/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, albumCtrl *controllers.AlbumController) {
	// Album routes
	albumRoutes := router.Group("/albums")
	{
		albumRoutes.GET("", albumCtrl.GetAlbums)          // GET /albums
		albumRoutes.GET("/:id", albumCtrl.GetAlbum)       // GET /albums/:id
		albumRoutes.POST("", albumCtrl.CreateAlbum)       // POST /albums
		albumRoutes.PUT("/:id", albumCtrl.UpdateAlbum)    // PUT /albums/:id
		albumRoutes.DELETE("/:id", albumCtrl.DeleteAlbum) // DELETE /albums/:id
	}
}
