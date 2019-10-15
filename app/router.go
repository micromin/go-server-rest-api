package app

import (
	"github.com/gin-gonic/gin"
	"gos/app/auth"
	"gos/app/controller"
	"gos/app/models"
	"net/http"
)

type Router struct {
	Engine     *gin.Engine
	Controller controller.IAppController
	Auth       auth.IAuth
}

// NewRouter creates a new router
func NewRouter(controller controller.IAppController, auth auth.IAuth) *Router {
	engine := gin.Default()

	router := &Router{
		Engine:     engine,
		Controller: controller,
		Auth:       auth,
	}

	engine.Use(addContentTypeHeader)
	engine.Use(handle404)

	registerRoutes(engine, router)

	return router
}

func addContentTypeHeader(context *gin.Context) {
	context.Header("Content-Type", "application/json; charset=utf-8")
}

func (r *Router) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader(auth.TokenHeader)
		_, err := r.Auth.AuthenticateUser(ctx, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &models.Response{
				Message: "UNAUTHORIZED ACCESS",
				Errors:  []string{err.Error()},
			})
			return
		}
	}
}

func handle404(ctx *gin.Context) {
	if ctx.Writer.Status() == http.StatusNotFound {
		ctx.AbortWithStatusJSON(http.StatusNotFound, &models.Response{
			Message: "page not found",
		})
		return
	}
}

func registerRoutes(engine *gin.Engine, router *Router) {
	api := engine.Group("/api/auth")
	{
		api.POST("/register", router.Controller.Register)
		api.POST("/login", router.Controller.Login)
	}

	// basic auth
	secured := engine.Group("/api/secured")
	{
		secured.Use(router.authMiddleware())
		{
			secured.POST("/tasks", router.Controller.AddTask)
			secured.GET("/tasks", router.Controller.GetTasks)
			secured.GET("/tasks/:taskId", router.Controller.GetTask)
		}
	}
}
