package api

import (
	"log"
	"net/http"

	"github.com/Masha003/Golang-gateway/internal/api/controllers"
	"github.com/Masha003/Golang-gateway/internal/api/middleware"
	"github.com/Masha003/Golang-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

func NewHttpServer(cfg config.Config, userController controllers.UserController) *http.Server {
	log.Print("Creating new server")

	e := gin.Default()
	r := e.Group("/api")

	// Register routes
	registerUserRoutes(r, cfg, userController)

	return &http.Server{
		Addr:    cfg.HttpPort,
		Handler: e,
	}
}

func registerUserRoutes(router *gin.RouterGroup, cfg config.Config, c controllers.UserController) {
	r := router.Group("/users")
	r.POST("/register", c.Register)
	r.POST("/login", c.Login)
	r.GET("/", c.GetAll)
	r.GET("/:id", c.GetById)

	pr := r.Use(middleware.JwtAuth(cfg.Secret))
	pr.DELETE("/:id", c.Delete)
}
