package controllers

import (
	"gin-freemarket-app/dto"
	"gin-freemarket-app/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IAuthController interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
}

func NewAuthController(service services.IAuthService) IAuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var input dto.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.Signup(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
