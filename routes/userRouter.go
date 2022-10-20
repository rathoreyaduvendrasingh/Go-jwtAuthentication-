package routes

import(
	controller "github.com/yaduvendra/E-commerce/controllers"
	middleware "github.com/yaduvendra/E-commerce/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/users/:userID",controller.GetUser())
}