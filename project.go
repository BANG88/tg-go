package main

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

/**
 * handleListProject list all projects as keyboard
 */
func (app *App) handleListProject(message *tgbotapi.Message) {
	innerJob, err := app.jenkins.GetAllJobNames()
	if err != nil {
		return
	}
	var data []string
	for _, job := range innerJob {
		data = append(data, fmt.Sprintf("/%s %s %s", app.commands.Project, app.operators.Build, job.Name))
	}
	var jobKeyboard = app.makeKeyboard(data, 3)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyMarkup = jobKeyboard
	app.bot.Send(msg)
}

/**
 * handleListProjectStatus list all projects's status as table
 */
func (app *App) handleListProjectStatus(message *tgbotapi.Message) {
	jobs, err := app.jenkins.GetAllJobs()
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
	app.bot.Send(msg)
}

/**
 * handle project
 */
func (app *App) handleProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	switch args {
	case app.operators.List:
		app.handleListProject(message)
	case "":
		app.handleListProjectStatus(message)
	}
}
