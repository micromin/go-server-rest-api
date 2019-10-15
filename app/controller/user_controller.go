package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gos/app/models"
	"gos/app/repo"
	"strconv"
	"time"
)

type IUserController interface {
	GetUsers(ctx *gin.Context)
	AddUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}

type UserController struct {
	userRepo repo.IUserRepo
}

func NewController(userRepo repo.IUserRepo) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

func (c *UserController) AddUser(ctx *gin.Context) {
	user := new(models.User)
	if err := ctx.BindJSON(user); err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("invalid request", err))
		return
	}

	user.DateUpdated = time.Now().Unix()
	user.DateCreated = time.Now().Unix()
	user.FailedLoginAttempt = 0
	user.LastLogin = 0

	_, err := c.userRepo.AddUser(ctx, *user)
	if err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("failed to add user", err))
		return
	}
	ctx.JSON(200, &models.Response{
		Message: "successfully added user",
	})
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	params := struct {
		LastId int64 `form:"lastId"`
		Limit  int   `form:"limit"`
	}{
		LastId: 0,
		Limit:  100,
	}

	err := ctx.BindQuery(&params)
	if err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("failed to extract query params", err))
		return
	}

	if params.Limit > 100 {
		params.Limit = 100
	} else if params.Limit <= 0 {
		params.Limit = 100
	}

	users, err := c.userRepo.GetUsers(ctx, params.LastId, params.Limit)
	if err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("failed to get users", err))
		return
	}

	lenUsers := len(users)

	var nextUserId int64
	if users != nil && lenUsers > 0 {
		nextUserId = users[lenUsers-1].UserId
	}

	ctx.JSON(200, &models.Response{
		Message: fmt.Sprintf("successfully retrieved %d user/s", lenUsers),
		Data: &models.Paged{
			Items:  users,
			LastId: nextUserId,
			Limit:  params.Limit,
		},
	})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	if len(userId) <= 0 {
		ctx.AbortWithStatusJSON(400, getErrorResponse("invalid request", errors.New("user id is missing")))
		return
	}

	userIdVal, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("invalid request", errors.New("user id is invalid")))
		return
	}

	user, err := c.userRepo.GetUser(ctx, userIdVal)
	if err != nil {
		ctx.AbortWithStatusJSON(400, getErrorResponse("failed to get user", err))
		return
	}

	ctx.JSON(200, &models.Response{
		Message: fmt.Sprintf("successfully retrieved user with Id %d", userIdVal),
		Data:    user,
	})
}

func getErrorResponse(msg string, err error) *models.Response {
	return &models.Response{
		Message: msg,
		Errors:  []string{err.Error()},
	}
}

var _ = (*UserController)(nil)
