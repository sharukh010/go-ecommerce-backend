package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sharukh010/go-ecommerce/controllers"
	middleware "github.com/sharukh010/go-ecommerce/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("/users/signup",controllers.SignUp())
	incomingRoutes.POST("/users/login",controllers.Login())
	incomingRoutes.GET("/users/search",controllers.SearchProductByQuery())
	incomingRoutes.GET("/users/productview",controllers.GetProducts())

	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.POST("/admin/addproduct",controllers.AddProduct())
}
