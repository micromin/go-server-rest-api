package controller

import "github.com/gin-gonic/gin"

type IUserController interface {
	GetUsers(ctx *gin.Context)
	AddUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}

type UserController struct {

}

func (c *UserController) AddUser(ctx *gin.Context)  {

}

func (c *UserController) GetUsers(ctx *gin.Context)  {

}

func (c *UserController) GetUser(ctx *gin.Context)  {

}

var _ = (*UserController)(nil)