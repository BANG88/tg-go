package main

import (
	"log"

	"github.com/bndr/gojenkins"
)

// App bot app
type App struct {
	conf Conf
}

// create jenkins instance
func (app *App) getJenkinsInstance() *gojenkins.Jenkins {
	var conf = app.conf
	jenkins, err := gojenkins.CreateJenkins(nil, conf.Jenkins.Server, conf.Jenkins.Username, conf.Jenkins.Password).Init()
	if err != nil {
		log.Fatalf("can not create jenkins instance %s", err)
	}
	return jenkins
}
func (app *App) start() {
	user := User{
		ID:      1,
		Name:    app.conf.SuperAdmin,
		IsAdmin: true,
	}
	app.init(user)
}

func main() {
	var app = App{
		conf: GetConf(),
	}
	app.start()
	db := getDbContext()
	defer db.Close()
	var user User
	db.One("Name", app.conf.SuperAdmin, &user)
	log.Printf("find user: %s", user.Name)
}
