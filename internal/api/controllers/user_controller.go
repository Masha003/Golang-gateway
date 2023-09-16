package controllers

import (
	"log"
	"net/http"

	"github.com/Masha003/Golang-gateway/internal/models"
	"github.com/Masha003/Golang-gateway/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewUserController(service service.UserService) UserController {
	log.Print("Creating new user controller")

	return &userController{
		service: service,
	}
}

type userController struct {
	service service.UserService
}

func (c *userController) GetAll(ctx *gin.Context) {
	query := models.PaginationQuery{}
	err := ctx.BindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := c.service.FindAll(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *userController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Register(ctx *gin.Context) {
	var user models.RegisterUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Register(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, token)
}

func (c *userController) Login(ctx *gin.Context) {
	var user models.LoginUser
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
