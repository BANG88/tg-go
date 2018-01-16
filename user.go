package main

import "log"

// User user
type User struct {
	ID      int
	Name    string `storm:"unique"`
	IsAdmin bool
}

func (app *App) init(user User) {
	db := getDbContext()
	defer db.Close()
	db.Init(&User{})

	err := db.Save(&user)
	if err != nil {
		log.Printf("save user failed: %s", err)
	}
}

func (user *User) findAll() {
	db := getDbContext()
	defer db.Close()
	// TODO:
}
