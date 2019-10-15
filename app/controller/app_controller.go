package controller

import (
	"github.com/gin-gonic/gin"
	"gos/app/auth"
	"gos/app/repo"
)

// IAppController is the main interface for app controllers
type IAppController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)

	GetTasks(ctx *gin.Context)
	GetTask(ctx *gin.Context)
	AddTask(ctx *gin.Context)
}

// AppController holds the repo connection and auth service
type AppController struct {
	appRepo repo.IAppRepo
	auth    auth.IAuth
}

// NewAppController returns a new controller for the app
func NewAppController(userRepo repo.IAppRepo, auth auth.IAuth) *AppController {
	return &AppController{
		appRepo: userRepo,
		auth:    auth,
	}
}
