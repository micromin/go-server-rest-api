package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gos/app/models"
)

type IAppRepo interface {
	AddUser(ctx context.Context, user models.User) (sql.Result, error)
	UpdateUser(ctx context.Context, user models.User) (sql.Result, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	GetAllTasks(ctx context.Context, lastTaskId int64, userId int64, limit int) ([]models.Task, error)
	GetTaskById(ctx context.Context, taskId int64, userId int64) (*models.Task, error)
	AddTask(ctx context.Context, task models.Task) (sql.Result, error)

	Close() error
}

type AppRepo struct {
	con               *sqlx.DB
	createUserStm     *sql.Stmt
	updateUserStm     *sql.Stmt
	getUserByEmailStm *sql.Stmt

	getAllTaskStm  *sql.Stmt
	getTaskByIdStm *sql.Stmt
	addTaskStm     *sql.Stmt
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

const createUserStatement = `insert into GOS_USER (name, email, password, last_login, failed_login_attempt, date_created, date_updated) VALUES (?, ?, ?, ?, ?, ?, ?)`
const updateUserStatement = `update GOS_USER set name = ?, email =?, password = ?, last_login = ?, failed_login_attempt = ? , date_created = ?, date_updated = ? where user_id = ?`
const getUserByEmailStatement = `select user_id, name, email, password, last_login, failed_login_attempt, date_created, date_updated from GOS_USER where email = ?`

const insertTaskStatement = `insert into GOS_TASK (user_id, title, description, date_created, date_updated, due_date, date_complete) VALUES (?, ?, ?, ?, ?, ?, ?)`
const getTasksStatement = `select task_id, user_id, title, description, date_created, date_updated, due_date, date_complete from GOS_TASK where task_id > ? and user_id = ? order by task_id desc limit ?`
const getTaskByIdStatement = `select task_id, user_id, title, description, date_created, date_updated, due_date, date_complete from GOS_TASK where task_id = ? and user_id = ?`

func NewAppRepo(dbConfig DbConfig) (*AppRepo, error) {
	name := dataStoreName(dbConfig)
	con, err := sqlx.Connect("mysql", name)
	if err != nil {
		return nil, err
	}

	createUserStm, err := con.Prepare(createUserStatement)
	if err != nil {
		return nil, err
	}

	updateUserStm, err := con.Prepare(updateUserStatement)
	if err != nil {
		return nil, err
	}

	getUserByEmailStm, err := con.Prepare(getUserByEmailStatement)
	if err != nil {
		return nil, err
	}

	getAllTaskStm, err := con.Prepare(getTasksStatement)
	if err != nil {
		return nil, err
	}

	getTaskByIdStm, err := con.Prepare(getTaskByIdStatement)
	if err != nil {
		return nil, err
	}

	addTaskStm, err := con.Prepare(insertTaskStatement)
	if err != nil {
		return nil, err
	}

	return &AppRepo{
		con:               con,
		createUserStm:     createUserStm,
		updateUserStm:     updateUserStm,
		getUserByEmailStm: getUserByEmailStm,
		getAllTaskStm:     getAllTaskStm,
		getTaskByIdStm:    getTaskByIdStm,
		addTaskStm:        addTaskStm,
	}, nil
}

func (r *AppRepo) AddUser(ctx context.Context, user models.User) (sql.Result, error) {
	return r.createUserStm.Exec(user.Name, user.Email, user.Password, user.LastLogin, user.FailedLoginAttempt, user.DateUpdated, user.DateUpdated)
}

func (r *AppRepo) UpdateUser(ctx context.Context, user models.User) (sql.Result, error) {
	return r.updateUserStm.Exec(user.Name, user.Email, user.Password, user.LastLogin, user.FailedLoginAttempt, user.DateUpdated, user.DateUpdated, user.UserId)
}

func (r *AppRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.getUserByEmailStm.QueryRow(email)

	user, err := scanRowUser(row)
	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("user not found")
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (r *AppRepo) GetAllTasks(ctx context.Context, lastTaskId int64, userId int64, limit int) ([]models.Task, error) {
	rows, err := r.getAllTaskStm.Query(lastTaskId, userId, limit)

	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, 0)
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		task, err := scanRowTask(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (r *AppRepo) GetTaskById(ctx context.Context, taskId int64, userId int64) (*models.Task, error) {
	row := r.getTaskByIdStm.QueryRow(taskId, userId)

	task, err := scanRowTask(row)
	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("task not found")
	case nil:
		return task, nil
	default:
		return nil, err
	}
}

func (r *AppRepo) AddTask(ctx context.Context, task models.Task) (sql.Result, error) {
	return r.addTaskStm.Exec(task.UserId, task.Title, task.Description, task.DateCreated, task.DateUpdated, task.DueDate, task.DateCompleted)
}

func (r *AppRepo) Close() error {
	return r.con.Close()
}

func scanRowUser(s RowScanner) (*models.User, error) {
	var (
		userId             int64
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

func scanRowTask(s RowScanner) (*models.Task, error) {
	var (
		taskId        int64
		userId        int64
		title         string
		description   string
		dateCreated   int64
		dateUpdated   int64
		dueDate       int64
		dateCompleted int64
	)
	if err := s.Scan(&taskId, &userId, &title, &description, &dateCreated, &dateUpdated, &dueDate, &dateCompleted); err != nil {
		return nil, err
	}

	return &models.Task{
		TaskId:        taskId,
		UserId:        userId,
		Title:         title,
		Description:   description,
		DateCreated:   dateCreated,
		DateUpdated:   dateUpdated,
		DueDate:       dueDate,
		DateCompleted: dateCompleted,
	}, nil
}

func dataStoreName(dbConfig DbConfig) string {
	return fmt.Sprintf("%s:%s@(%s:%v)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DatabaseName)
}

var _ = (*AppRepo)(nil)
