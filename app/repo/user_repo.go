package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gos/app/models"
)

type IUserRepo interface {
	AddUser(ctx context.Context, user models.User) (sql.Result, error)
	GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error)
	GetUser(ctx context.Context, userId int64) (*models.User, error)
	Close() error
}

type UserRepo struct {
	con            *sqlx.DB
	saveUserStm    *sql.Stmt
	getUserStm     *sql.Stmt
	getOneUsersStm *sql.Stmt
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type DbConfig struct {
	Host         string `required:"true"`
	Port         int    `required:"true"`
	DatabaseName string `required:"true"`
	User         string `required:"true"`
	Password     string `required:"true"`
}

const insertUserStatement = `insert into GOS_USER (name, email, password, last_login, failed_login_attempt, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?, ?)`
const getUsersStatement = `select user_id, name, email, password, last_login, failed_login_attempt, date_created, date_updated from GOS_USER where user_id > ? order by user_id desc limit ?`
const getOneUserStatement = `select user_id, name, email, password, last_login, failed_login_attempt, date_created, date_updated from GOS_USER where user_id = ?`

func NewUserRepo(dbConfig DbConfig) (*UserRepo, error) {
	name := dataStoreName(dbConfig)
	con, err := sqlx.Connect("mysql", name)
	if err != nil {
		return nil, err
	}

	saveUserStm, err := con.Prepare(insertUserStatement)
	if err != nil {
		return nil, err
	}

	getUsersStm, err := con.Prepare(getUsersStatement)
	if err != nil {
		return nil, err
	}

	getOneUsersStm, err := con.Prepare(getOneUserStatement)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		con:            con,
		getUserStm:     getUsersStm,
		saveUserStm:    saveUserStm,
		getOneUsersStm: getOneUsersStm,
	}, nil
}
func (r *UserRepo) AddUser(ctx context.Context, user models.User) (sql.Result, error) {
	return r.saveUserStm.Exec(user.Name, user.Email, user.Password, user.LastLogin, user.FailedLoginAttempt, user.DateUpdated,user.DateUpdated)
}

func (r *UserRepo) GetUsers(ctx context.Context, lastUserId int64, limit int) ([]models.User, error) {
	rows, err := r.getUserStm.Query(lastUserId, limit)

	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0)
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		user, err := scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}
		users = append(users, *user)
	}

	return users, nil
}

func (r *UserRepo) GetUser(ctx context.Context, userId int64) (*models.User, error) {
	row := r.getOneUsersStm.QueryRow(userId)

	user, err := scanRow(row)
	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("user not found")
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (r *UserRepo) Close() error {
	return r.con.Close()
}

func scanRow(s RowScanner) (*models.User, error) {
	var (
		userId                 int64
		name               string
		email              string
		password           string
		lastLogin          int
		failedLoginAttempt int
		dateCreated        int64
		dateUpdated        int64
	)
	if err := s.Scan(&userId, &name, &email, &password, &lastLogin, &failedLoginAttempt, &dateCreated, &dateUpdated); err != nil {
		return nil, err
	}

	return &models.User{
		UserId:             userId,
		Name:               email,
		Email:              email,
		Password:           password,
		LastLogin:          lastLogin,
		FailedLoginAttempt: failedLoginAttempt,
		DateCreated:        dateCreated,
		DateUpdated:        dateUpdated,
	}, nil
}

func dataStoreName(dbConfig DbConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%v)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName)
}

var _ = (*UserRepo)(nil)
