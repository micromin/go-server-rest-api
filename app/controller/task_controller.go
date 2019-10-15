package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gos/app/auth"
	"gos/app/models"
	"net/http"
	"strconv"
	"time"
)

// swagger:operation GET /api/secured/tasks GetTasks
//
// GetTasks gets tasks for the logged in user
// ---
// produces:
// - application/json
// parameters:
// - name: x-access-token
//   in: header
//   description: the access token
//   type: string
// - name: lastId
//   in: query
//   description: the id of the last task in the response
//   type: string
// - name: limit
//   in: query
//   description: the page size value
//   type: integer
// responses:
//  '200':
//    description: successful operation
//    schema:
//     $ref: '#/definitions/Response'
//  '400':
//    description: invalid request
//    schema:
//     $ref: '#/definitions/Response'
//  '500':
//    description: internal server error
//    schema:
//     $ref: '#/definitions/Response'
//  '401':
//    description: unauthorized access
//    schema:
//     $ref: '#/definitions/Response'
func (c *AppController) GetTasks(ctx *gin.Context) {
	params := struct {
		LastId int64 `form:"lastId"`
		Limit  int   `form:"limit"`
	}{
		LastId: 0,
		Limit:  100,
	}

	err := ctx.BindQuery(&params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to extract query params", err))
		return
	}

	if params.Limit > 100 {
		params.Limit = 100
	} else if params.Limit <= 0 {
		params.Limit = 100
	}

	claims := ctx.MustGet("claims")
	claimsObj := claims.(*auth.Claims)

	tasks, err := c.appRepo.GetAllTasks(ctx, params.LastId, claimsObj.UserId, params.Limit)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to get tasks", err))
		return
	}

	lenTasks := len(tasks)

	var nextTaskId int64
	if tasks != nil && lenTasks > 0 {
		nextTaskId = tasks[lenTasks-1].UserId
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Message: fmt.Sprintf("successfully retrieved %d task/s", lenTasks),
		Data: &models.Paged{
			Items:  tasks,
			LastId: nextTaskId,
			Limit:  params.Limit,
		},
	})
}

// swagger:operation GET /api/secured/tasks/:taskId GetTask
//
// GetTask gets a task for the logged in user
// ---
// produces:
// - application/json
// parameters:
// - name: x-access-token
//   in: header
//   description: the access token
//   type: string
// responses:
//  '200':
//    description: successful operation
//    schema:
//     $ref: '#/definitions/Response'
//  '400':
//    description: invalid request
//    schema:
//     $ref: '#/definitions/Response'
//  '500':
//    description: internal server error
//    schema:
//     $ref: '#/definitions/Response'
//  '401':
//    description: unauthorized access
//    schema:
//     $ref: '#/definitions/Response'
func (c *AppController) GetTask(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	taskIdVal, err := strconv.ParseInt(taskId, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", errors.New("task id is invalid")))
		return
	}

	claims := ctx.MustGet("claims")
	claimsObj := claims.(*auth.Claims)

	task, err := c.appRepo.GetTaskById(ctx, taskIdVal, claimsObj.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to get task", err))
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Message: fmt.Sprintf("successfully retrieved task with id %d", taskIdVal),
		Data:    task,
	})
}

// swagger:operation POST /api/secured/tasks AddTask
//
// AddTask adds a task
// ---
// produces:
// - application/json
// parameters:
// - name: x-access-token
//   in: header
//   description: the access token
//   type: string
// - name: body
//   in: body
//   description: the task to be added
//   schema:
//    $ref: '#/definitions/Task'
// responses:
//  '201':
//    description: successful operation
//    schema:
//     $ref: '#/definitions/Response'
//  '400':
//    description: invalid request
//    schema:
//     $ref: '#/definitions/Response'
//  '500':
//    description: internal server error
//    schema:
//     $ref: '#/definitions/Response'
//  '401':
//    description: unauthorized access
//    schema:
//     $ref: '#/definitions/Response'
func (c *AppController) AddTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims")
	claimsObj := claims.(*auth.Claims)
	task := new(models.Task)
	if err := ctx.BindJSON(task); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", err))
		return
	}

	task.DateCreated = time.Now().Unix()
	task.DateUpdated = time.Now().Unix()
	task.UserId = claimsObj.UserId

	if len(task.Title) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("invalid request", errors.New("field title is required")))
		return
	}

	_, err := c.appRepo.AddTask(ctx, *task)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, getErrorResponse("failed to add task", err))
		return
	}

	ctx.JSON(http.StatusCreated, &models.Response{
		Message: "successfully added a task",
	})
}

func getErrorResponse(msg string, err error) *models.Response {
	if err == nil {
		return &models.Response{
			Message: msg,
		}
	}

	return &models.Response{
		Message: msg,
		Errors:  []string{err.Error()},
	}
}

var _ = (*AppController)(nil)
