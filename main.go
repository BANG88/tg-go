package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bang88/gojenkins"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Commands supported commands
type Commands struct {
	Admin   string
	Project string
	Help    string
}

// Operators Operators
type Operators struct {
	List   string
	Add    string
	Remove string
	Build  string
}

// App bot app
type App struct {
	conf      Conf
	commands  Commands
	operators Operators
	bot       *tgbotapi.BotAPI
	jenkins   *gojenkins.Jenkins
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

func (app *App) handleHelp(message *tgbotapi.Message) {
	var args = message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(`
Usage:
- admin ls (List all administrators)
- help (Get help)
- project ls (List all projects)
			`))
		// msg.ReplyToMessageID = message.MessageID
		app.bot.Send(msg)
	}
}
func (app *App) handleAdmin(message *tgbotapi.Message) {
	var args = message.CommandArguments()
	if args != "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(`command: %s,args: %s,`, app.commands.Admin, args))
		msg.ReplyToMessageID = message.MessageID
		app.bot.Send(msg)
	}
}

/**
 * getCommandArguments getCommandArguments
 */
func (app *App) getCommandArguments(message *tgbotapi.Message) string {
	args := strings.TrimSpace(message.CommandArguments())
	return args
}

// start bot routine
func (app *App) startBot() {
	bot, err := tgbotapi.NewBotAPI(app.conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	// app bot
	app.bot = bot
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		user := app.findUser(update.Message.Chat.UserName)
		if !user.IsAdmin {
			continue
		}
		// handle command
		command := strings.ToLower(update.Message.Command())
		switch command {
		case app.commands.Admin:
			app.handleAdmin(update.Message)
		case app.commands.Help:
			app.handleHelp(update.Message)
		case app.commands.Project:
			app.handleProject(update.Message)
		}
	}
}

// runnnn
func main() {
	var commands = Commands{
		Admin:   "admin",
		Help:    "help",
		Project: "project",
	}
	var operators = Operators{
		List:   "ls",
		Add:    "add",
		Remove: "rm",
		Build:  "build",
	}
	var app = App{
		conf:      GetConf(),
		commands:  commands,
		operators: operators,
	}
	app.jenkins = app.getJenkinsInstance()
	app.start()
	app.startBot()
}
