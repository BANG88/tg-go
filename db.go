package main

import (
	"log"

	"github.com/asdine/storm"
)

// getDbContext
func getDbContext() *storm.DB {
	var conf = GetConf()
	db, err := storm.Open(conf.DbPath, storm.BoltOptions(0600, nil))
	if err != nil {
		log.Printf("can not open database %s", err)
	}
	return db
}
