package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/identityprovider"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
)

type UserController struct {
	identityproviderAdapter identityprovider.CognitoAdapter
}

func NewUserController(identityproviderAdapter identityprovider.CognitoAdapter) *UserController {
	return &UserController{
		identityproviderAdapter: identityproviderAdapter,
	}
}

// @Summary Sign up a new user
// @Description Register a new user with email and password
// @Accept  json
// @Produce  json
// @Param body body entity.UserRegistration true "User registration info"
// @Success  200 {object} map[string]string{"message": "User registered successfully"} "Successfully registered"
// @Failure  400 {object} map[string]string{"error": "Bad request"} "Invalid input"
// @Failure  500 {object} map[string]string{"error": "Internal server error"} "Internal server error"
// @Router /signup [post]
func (c *UserController) SignUp(ctx *gin.Context) {
	var userRegistration entity.UserRegistration
	if err := ctx.ShouldBindJSON(&userRegistration); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.identityproviderAdapter.Register(ctx, userRegistration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary Log in with email and password
// @Description Authenticate user with email and password
// @Accept  json
// @Produce  json
// @Param body body entity.UserLogin true "User login info"
// @Success  200 {object} map[string]string{"Authorization": "JWT"} "Successfully logged in"
// @Failure  400 {object} map[string]string{"error": "Bad request"} "Invalid input"
// @Failure  500 {object} map[string]string{"error": "Internal server error"} "Internal server error"
// @Router /login [post]
func (c *UserController) LogIn(ctx *gin.Context) {
	var userLogin entity.UserLogin
	fmt.Println(userLogin.Email)
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error here": err.Error()})
		return
	}

	result, err := c.identityproviderAdapter.Login(ctx, userLogin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Authorization": result})
}
