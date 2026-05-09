package router

import (
	// "net/http"

	"github.com/bigthamm/task-his/internal/delivery/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Staff handler.StaffHandler
}

// func NewRouter(h Handlers) *gin.Engine {
// 	routers := gin.New()

// 	apiStaffGroup := routers.Group("/api")
// 	{
// 		apiStaffGroup.POST("/create", h.Staff.Create)
// 	}

// 	return routers
// }

func SetStaffRoute(router *gin.Engine, handler *handler.StaffHandler) {
	apiStaffGroup := router.Group("/api/v1")
	{
		apiStaffGroup.POST("/create-staff", handler.Create)
		apiStaffGroup.POST("/staff-login", handler.Login)
	}
}
