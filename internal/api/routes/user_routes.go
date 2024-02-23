package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/interfaces"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/api/middleware"
	"github.com/Zeta-Manu/Backend/internal/config"
)

func InitUserRoutes(router *gin.Engine, cfg config.AppConfig, idpAdapter interfaces.IIdentityProvider) {
	idpController := controllers.NewUserController(idpAdapter)

	// TODO: Fixed Swagger
	user := router.Group("/api/user")
	{
		user.POST("/signup", idpController.SignUp)
		user.POST("/login", idpController.LogIn)
		user.POST("/confirm", idpController.ConfirmRegistration)
		user.POST("/resend-confirm", idpController.ResendConfirmationCode)
		user.GET("/:email", idpController.GetUser)
		user.POST("/forgot-password", idpController.ForgotPassword)
		user.POST("/confirm-forgot-password", idpController.ConfirmForgotPassword)
		// route with middleware
		user.POST("/change-password", middleware.AuthenticationMiddleware(cfg), idpController.ChangePassword)
	}
}
