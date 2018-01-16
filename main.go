package main

import (
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/bndr/gojenkins"
	bolt "github.com/coreos/bbolt"
)

// App bot app
type App struct {
	conf Conf
}

// getDbContext
func (app *App) getDbContext() *storm.DB {
	var conf = app.conf
	db, err := storm.Open(conf.DbPath, storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	if err != nil {
		log.Printf("can not open database %s", err)
	}
	return db
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
	// db := app.getDbContext()
	// defer db.Close()
	// var user User
	// db.One("Name", app.conf.SuperAdmin, &user)
	// log.Printf("find user: %s", user.Name)
}
