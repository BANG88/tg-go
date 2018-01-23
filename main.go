package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olekukonko/tablewriter"
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
func (app *App) makeTable(data [][]string) string {
	file, err := ioutil.TempFile(os.TempDir(), "bot-")
	if err != nil {
		log.Printf("err: %s", err)
		return ""
	}

	table := tablewriter.NewWriter(file)
	table.SetHeader([]string{"#", "Project", "Latest Build"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()         // Send output

	buf := bytes.NewBuffer(nil)
	f, _ := os.Open(file.Name())
	io.Copy(buf, f)
	f.Close()
	file.Close()
	defer os.Remove(file.Name())
	return fmt.Sprintf("`%s`", string(buf.Bytes()))
}
func (app *App) handleHelp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var args = message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(`
Usage:
- admin ls (List all administrators)
- help (Get help)
- project ls (List all projects)
			`))
		// msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	}
}
func (app *App) handleAdmin(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var args = message.CommandArguments()
	if args != "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(`
			command: %s,
			args: %s,
			`, app.commands.Admin, args))
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	}
}
func (app *App) handleProject(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	j := app.getJenkinsInstance()
	var args = message.CommandArguments()
	if args == "" {
		jobs, err := j.GetAllJobs()
		if err != nil {
			return
		}
		var data [][]string
		for index, job := range jobs {
			build, err := job.GetLastBuild()
			result := "N/A"
			if err != nil {
				log.Printf("error: %s", err)
			} else {
				result = build.GetResult()
			}
			id := fmt.Sprintf("%v", index+1)
			name := fmt.Sprintf("%s", job.GetName())
			// name := fmt.Sprintf("[%s](http://www.example.com/)", job.GetName())
			data = append(data, []string{id, name, result})
		}
		var str = app.makeTable(data)
		log.Printf("str: %s", str)
		msg := tgbotapi.NewMessage(message.Chat.ID, str)
		msg.ParseMode = tgbotapi.ModeMarkdown
		// msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	} else {
		innerJob, err := j.GetAllJobNames()
		if err != nil {
			return
		}
		var keyboardButton []tgbotapi.KeyboardButton
		for _, job := range innerJob {
			keyboardButton = append(keyboardButton, tgbotapi.NewKeyboardButton(fmt.Sprintf("/%s %s %s", app.commands.Project, app.operators.Build, job.Name)))
		}
		var keyboard [][]tgbotapi.KeyboardButton
		keyboard = append(keyboard, keyboardButton)
		var jobKeyboard = tgbotapi.ReplyKeyboardMarkup{
			Keyboard:        keyboard,
			OneTimeKeyboard: true,
			Selective:       true,
			ResizeKeyboard:  true,
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
		msg.ReplyMarkup = jobKeyboard
		bot.Send(msg)
	}
}

// start bot routine
func (app *App) startBot() {
	bot, err := tgbotapi.NewBotAPI(app.conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

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
		user := app.findUser(update.Message.Chat.UserName)
		if !user.IsAdmin {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		// handle command
		command := strings.ToLower(update.Message.Command())
		switch command {
		case app.commands.Admin:
			app.handleAdmin(bot, update.Message)
		case app.commands.Help:
			app.handleHelp(bot, update.Message)
		case app.commands.Project:
			app.handleProject(bot, update.Message)
		}
	}
}
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
	app.start()
	app.startBot()
}
