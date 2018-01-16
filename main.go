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
	var conf = app.conf.getConf()
	jenkins, err := gojenkins.CreateJenkins(nil, conf.Jenkins.Server, conf.Jenkins.Username, conf.Jenkins.Password).Init()
	if err != nil {
		log.Fatalf("can not create jenkins instance %s", err)
	}
	return jenkins
}
func (app *App) start() {
}

func main() {
	log.Printf("hello bot")
}
