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
func (app *App) findUser(username string) *User {
	db := getDbContext()
	defer db.Close()
	var user User
	db.One("Name", username, &user)

	return &user
}
func (user *User) findAll() {
	db := getDbContext()
	defer db.Close()
	// TODO:
}
