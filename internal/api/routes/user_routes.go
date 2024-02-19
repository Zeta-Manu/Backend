package routes

import (
	"github.com/gin-gonic/gin"

	docs "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/interfaces"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func InitUserRoutes(router *gin.Engine, idpAdapter interfaces.IIdentityProvider) {
	idpController := controllers.NewUserController(idpAdapter)

	// TODO: Fixed Swagger
	docs.SwaggerInfo.BasePath = "/user"
	user := router.Group("/user")
	{
		user.POST("/signup", idpController.SignUp)
		user.POST("/login", idpController.LogIn)
	}
}
