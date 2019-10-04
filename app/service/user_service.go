package controller

import "github.com/gin-gonic/gin"

type IUserService interface {
	GetUsers(ctx *gin.Context)
	AddUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}

type UserService struct {

}

func (s *UserService) AddUser(ctx *gin.Context)  {

}

func (s *UserService) GetUsers(ctx *gin.Context)  {

}

func (s *UserService) GetUser(ctx *gin.Context)  {

}

var _ = (*UserService)(nil)