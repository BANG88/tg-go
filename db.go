package main

import (
	"log"
	"time"

	"github.com/asdine/storm"
	bolt "github.com/coreos/bbolt"
)

// getDbContext
func getDbContext() *storm.DB {
	var conf = GetConf()
	db, err := storm.Open(conf.DbPath, storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	if err != nil {
		log.Printf("can not open database %s", err)
	}
	return db
}
