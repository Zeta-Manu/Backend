package routes

import (
	"github.com/gin-gonic/gin"
  swaggerfiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"

  docs "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/config"
)

func InitUserRoute(router *gin.Engine, controller *controllers.UserController) {
  docs.SwaggerInfo.BasePath = "/user"
  user := router.Group("/user")
  {
	  user.POST("/signup", controller.SignUp)
	  user.POST("/login", controller.LogIn)
  }

  url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
}
