// GO Simple Server (GOS)
//
// Go simple (gos) is a simple GO server for adding todos.
//
//     Schemes: http
//     Host: localhost:8080
//     Version: 0.1.0
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gos/app"
	"gos/app/auth"
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
	userRepo, err := repo.NewAppRepo(repo.DbConfig{
		Host:         "127.0.0.1", // get the host from env variable
		Port:         3306,
		DatabaseName: "gos",
		User:         "gos",
		Password:     "1234", // do not put password in the code, get it from env variable
	})

	if err != nil {
		die(err)
	}

	authService := auth.NewAuth(userRepo, "some-key")
	appController := controller.NewAppController(userRepo, authService)
	router := app.NewRouter(appController, authService)

	err = router.Engine.Run(":8080")
	if err != nil {
		die(err)
	}
}
