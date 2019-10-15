package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gos/app"
	"gos/app/controller"
	"gos/app/repo"
	"os"
)


func die(err error) {
	if err == nil {
		err = errors.New("die with no error..")
	}

	fmt.Printf(err.Error())

	os.Exit(1)
}

func main() {
	userRepo, err := repo.NewUserRepo(repo.DbConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		DatabaseName: "test_go",
		User:     "demo_go",
		Password: "1234",
	})

	if err != nil {
		die(err)
	}

	appController := controller.NewController(userRepo)
	router := app.NewRouter(appController)

	err = router.Engine.Run(":8080")
	if err != nil {
		die(err)
	}
}
