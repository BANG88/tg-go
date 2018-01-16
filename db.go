package main

import (
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
)

func (app *App) getDbContext() *bolt.DB {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(app.conf.getConf().DbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	return db
}
func (app *App) findAllUser() {
	db := app.getDbContext()
	log.Printf("db %v", db)
}
