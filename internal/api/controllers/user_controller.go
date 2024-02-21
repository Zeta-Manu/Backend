package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/interfaces"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
)

type UserController struct {
	identityproviderAdapter interfaces.IIdentityProvider
}

func NewUserController(identityproviderAdapter interfaces.IIdentityProvider) *UserController {
	return &UserController{
		identityproviderAdapter: identityproviderAdapter,
	}
}

// @Summary		Sign up a new user
// @Description	Register a new user with email and password
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserRegistration											true	"User registration info"
// @Router			/signup [post]
func (c *UserController) SignUp(ctx *gin.Context) {
	var userRegistration entity.UserRegistration
	if err := ctx.ShouldBindJSON(&userRegistration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.identityproviderAdapter.Register(ctx, userRegistration)
	if err != nil {
		// FIX: Handle Additonal Error like InvalidPasswordException and other
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary		Log in with email and password
// @Description	Authenticate user with email and password
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserLogin									true	"User login info"
// @Router			/login [post]
func (c *UserController) LogIn(ctx *gin.Context) {
	var userLogin entity.UserLogin
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		// FIX: Handle Additional Error like NotAuthorizedException
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.identityproviderAdapter.Login(ctx, userLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", result)
	ctx.Status(http.StatusOK)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.identityproviderAdapter.GetUser(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UserController) ConfirmRegistration(ctx *gin.Context) {
	var userRegistrationConfirm entity.UserRegistrationConfirm
	if err := ctx.ShouldBindJSON(&userRegistrationConfirm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.identityproviderAdapter.ConfirmRegistration(ctx, userRegistrationConfirm)
	if err != nil {
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) ResendConfirmationCode(ctx *gin.Context) {
	var email entity.Email
	result, err := c.identityproviderAdapter.ResendConfirmationCode(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var email entity.Email
	if err := ctx.ShouldBindJSON(&email); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.identityproviderAdapter.ForgotPassword(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (c *UserController) ConfirmForgotPassword(ctx *gin.Context) {
	var userResetPassword entity.UserResetPassword
	if err := ctx.ShouldBindJSON(&userResetPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.identityproviderAdapter.ConfirmForgotPassword(ctx, userResetPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	var userChangePassword entity.UserChangePassword
	if err := ctx.ShouldBindJSON(&userChangePassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idToken := ctx.GetHeader("Authorization")

	err := c.identityproviderAdapter.ChangePassword(ctx, idToken, userChangePassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
