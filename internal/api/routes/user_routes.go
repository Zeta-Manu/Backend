package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/interfaces"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func InitUserRoute(router *gin.Engine, idpAdapter interfaces.IIdentityProvider) {
	idpController := controllers.NewUserController(idpAdapter)

	docs.SwaggerInfo.BasePath = "/user"
	user := router.Group("/user")
	{
		user.POST("/signup", idpController.SignUp)
		user.POST("/login", idpController.LogIn)
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
}
