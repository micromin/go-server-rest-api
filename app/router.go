package app

import (
	"github.com/gin-gonic/gin"
	"gos/app/controller"
)

type Router struct {
	Engine     *gin.Engine
	Controller controller.IUserController
}

func NewRouter(controller controller.IUserController) *Router {
	engine := gin.Default()

	router := &Router{
		Engine:     engine,
		Controller: controller,
	}

	engine.Use(addContentTypeHeader)

	registerRoutes(engine, router)

	return router
}

func addContentTypeHeader(context *gin.Context) {
	context.Header("Content-Type", "application/json; charset=utf-8")
}

func registerRoutes(engine *gin.Engine, router *Router) {
	engine.POST("/users", router.Controller.AddUser)
	engine.GET("/users", router.Controller.GetUsers)
	engine.GET("/users/:userId", router.Controller.GetUser)
}
