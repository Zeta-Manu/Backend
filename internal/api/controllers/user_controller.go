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
// @Tags User
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserRegistration											true	"User registration info"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router			/user/signup [post]
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

	response := gin.H{
		"data": result,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Log in with email and password
// @Description	Authenticate user with email and password
// @Tags User
// @Accept			json
// @Produce		json
// @Param			body	body		entity.UserLogin									true	"User login info"
// @Success 200 {object} entity.ResponseWrapper{data=entity.LoginResult}
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router			/user/login [post]
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

	response := gin.H{
		"data": result,
	}

	ctx.JSON(http.StatusOK, response)
}

// NOTE: Admin Thing
func (c *UserController) GetUser(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.identityproviderAdapter.GetUser(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary Confirm user registration
// @Description Confirm a user's registration using the provided confirmation information
// @Tags User
// @Accept json
// @Produce json
// @Param body body entity.UserRegistrationConfirm true "User registration confirmation info"
// @Success 200
// @Failure 400
// @Failure 408
// @Router /user/confirm [post]
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

// @Summary Resend confirmation code
// @Description Resend the confirmation code to the provided email address
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.Email true "Email address to resend the confirmation code to"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /user/resend-confirm [post]
func (c *UserController) ResendConfirmationCode(ctx *gin.Context) {
	var email entity.Email
	result, err := c.identityproviderAdapter.ResendConfirmationCode(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data": result,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Forgot Password
// @Description Initiate the password reset process for a user by sending a reset link to their email
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.Email true "Email address of the user"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /user/forgot-password [post]
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

	response := gin.H{
		"result": result,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Forgot Password
// @Description Initiate the password reset process for a user by sending a reset link to their email
// @Tags User
// @Accept json
// @Produce json
// @Param email body entity.UserResetPassword true "Email address of the user"
// @Success 200 {object} entity.ResponseWrapper
// @Failure 400 {object} entity.ErrorWrapper
// @Failure 500 {object} entity.ErrorWrapper
// @Router /user/confirm-forgot-password [post]
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

	response := gin.H{
		"data": result.String(),
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param body body entity.UserChangePassword true "User change password info"
// @Success  200
// @Failure  400 {object} entity.ErrorWrapper
// @Failure  401 {object} entity.ErrorWrapper
// @Failure  500 {object} entity.ErrorWrapper
// @Security BearerAuth
// @Router /user/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var userChangePassword entity.UserChangePassword
	if err := ctx.ShouldBindJSON(&userChangePassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	err := c.identityproviderAdapter.ChangePassword(ctx, token.(string), userChangePassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
