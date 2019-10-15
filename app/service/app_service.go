package service

import (
	"context"
	"database/sql"
	"gos/app/models"
	"gos/app/repo"
)

// IAppService is the main service for the app
type IAppService interface {
	AddUser(ctx context.Context, user models.User) (sql.Result, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	GetAllTasks(ctx context.Context, lastTaskId int64, userId int64, limit int) ([]models.Task, error)
	GetTaskById(ctx context.Context, taskId int64, userId int64) (*models.Task, error)
}

// AppService holds the repo
type AppService struct {
	appRepo *repo.AppRepo
}

// AddUser adds a user
func (s *AppService) AddUser(ctx context.Context, user models.User) (sql.Result, error) {
	return s.appRepo.AddUser(ctx, user)
}

// GetUserByEmail gets a user by email
func (s *AppService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.appRepo.GetUserByEmail(ctx, email)
}

// GetAllTasks gets all tasks for a user
func (s *AppService) GetAllTasks(ctx context.Context, lastTaskId int64, userId int64, limit int) ([]models.Task, error) {
	return s.appRepo.GetAllTasks(ctx, lastTaskId, userId, limit)
}

// GetTaskById gets a task by id
func (s *AppService) GetTaskById(ctx context.Context, taskId int64, userId int64) (*models.Task, error) {
	return s.appRepo.GetTaskById(ctx, taskId, userId)
}

// makes sure the interface is implemented by service
var _ = (*AppService)(nil)
