package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// User user
type User struct {
	ID      int    `storm:"id,increment"`
	Name    string `storm:"unique"`
	IsAdmin bool
}

func (app *App) init(user User) {
	app.addUser(user)
}

// addUser Add user
func (app *App) addUser(user User) error {

	u := app.findUser(user.Name)
	if &u != nil {
		return nil
	}
	db := getDbContext()
	defer db.Close()
	err := db.Save(&user)
	if err != nil {
		log.Printf("add user failed: %s", err)
		return err
	}
	return nil
}
func (app *App) findUser(username string) User {
	db := getDbContext()
	defer db.Close()
	var user User
	db.One("Name", username, &user)

	return user
}
func (app *App) removeUser(username string) error {
	users := app.findAllUsers()
	if len(users) == 1 {
		return errors.New("We need At least one admin")
	}
	user := app.findUser(username)
	if &user != nil {
		db := getDbContext()
		defer db.Close()
		err := db.DeleteStruct(&user)
		return err
	}
	return errors.New("Not found")
}
func (app *App) findAllUsers() []*User {
	db := getDbContext()
	defer db.Close()
	var users []*User
	err := db.All(&users)
	if err != nil {
		return nil
	}
	return users
}

// handleAdmin handleAdmin
func (app *App) handleAdmin(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	isAdd := strings.HasPrefix(args, app.operators.Add)
	isRemove := strings.HasPrefix(args, app.operators.Remove)
	switch true {
	// add admin
	case isAdd:
		username := app.getArgument(args)
		user := User{
			Name:    username,
			IsAdmin: true,
		}
		err := app.addUser(user)
		str := fmt.Sprintf("%s created", username)
		if err != nil {
			str = fmt.Sprintf("add user %s failed: %s", username, err)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, str)
		app.bot.Send(msg)
		// remove admin
	case isRemove:
		username := app.getArgument(args)
		if username == message.Chat.UserName {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Don't do this. give yourself a chance 😆")
			app.bot.Send(msg)
			return
		}
		err := app.removeUser(username)
		str := fmt.Sprintf("%s deleted", username)
		if err != nil {
			str = fmt.Sprintf("delete user %s failed: %s", username, err)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, str)
		app.bot.Send(msg)
	default:
		users := app.findAllUsers()
		if users != nil {
			var data [][]string
			for index, user := range users {
				data = append(data, []string{fmt.Sprintf("%v", index), user.Name})
			}
			var str = app.makeTable(data, []string{"#", "Name"})
			msg := tgbotapi.NewMessage(message.Chat.ID, str)
			msg.ParseMode = tgbotapi.ModeMarkdown
			app.bot.Send(msg)
		}
	}
}
