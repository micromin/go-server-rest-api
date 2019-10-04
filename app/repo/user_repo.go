package controller

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gos/app/models"
)

type IUserRepo interface {
	GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error)
	AddUserAddUser(ctx context.Context, user models.User) (sql.Result, error)
	GetUser(ctx context.Context, lastUserId int64, limit int) (*models.User, error)
}

type UserRepo struct {
	con         *sqlx.DB
	saveFeedStm *sql.Stmt
	getFeedsStm *sql.Stmt
}

type DbConfig struct {
	Host         string `required:"true"`
	Port         int    `required:"true"`
	DatabaseName string `required:"true"`
	User         string `required:"true"`
	Password     string `required:"true"`
}

const insertUserStatement = `insert into gos_users (name, email, password, ip, last_login, failed_login_attempt, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
const getUserStatement = `select * from gos_users where id > ? order by id desc limit ?`

func NewUserRepo(dbConfig DbConfig) (*UserRepo, error) {
	con, err := sqlx.Connect("mysql", dataStoreName(dbConfig))
	if err != nil {
		return nil, err
	}

	saveUserStm, err := con.Prepare(insertUserStatement)
	if err != nil {
		return nil, err
	}

	getUsersStm, err := con.Prepare(getUserStatement)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		con:         con,
		saveFeedStm: getUsersStm,
		getFeedsStm: saveUserStm,
	}, nil
}

func (s *UserRepo) AddUser(ctx context.Context, user models.User) (sql.Result, error) {
	return nil, nil
}

func (s *UserRepo) GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error) {
	return nil, nil
}

func (s *UserRepo) GetUser(ctx context.Context, lastUserId int64, limit int) (*models.User, error) {
	return nil, nil
}

func dataStoreName(dbConfig DbConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%v)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName)
}

var _ = (*UserRepo)(nil)
