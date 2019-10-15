package controller

import (
	"context"
	"database/sql"
	"gos/app/models"
	"gos/app/repo"
)

type IUserService interface {
	AddUser(ctx context.Context, user models.User) (sql.Result, error)
	GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error)
	GetUser(ctx context.Context, userId int64) (*models.User, error)
}

type UserService struct {
	userRepo *repo.UserRepo
}

func (s *UserService) AddUser(ctx context.Context, user models.User) (sql.Result, error) {
	return s.userRepo.AddUser(ctx, user)
}

func (s *UserService) GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error) {
	return s.userRepo.GetUsers(ctx, lastUserId, limit)
}

func (s *UserService) GetUser(ctx context.Context, userId int64) (*models.User, error) {
	return s.userRepo.GetUser(ctx, userId)
}

var _ = (*UserService)(nil)
