package main

import (
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
		Name:    app.conf.SuperAdmin,
		IsAdmin: true,
	}
	app.init(user)
}

func (app *App) handleHelp(message *tgbotapi.Message) {
	var data = [][]string{
		[]string{"/" + app.commands.Help, "🆘"},
		[]string{"/" + app.commands.Admin, "List all administrators"},
		[]string{"/" + app.commands.Project, "List all projects"},
	}

	var str = app.makeTable(data, []string{"Command", "Desc"})
	msg := tgbotapi.NewMessage(message.Chat.ID, str)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyToMessageID = message.MessageID
	app.bot.Send(msg)
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
		case "start":
			app.handleHelp(update.Message)
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
