package main

import (
	"fmt"
	"log"
	"regexp"

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
	var jobKeyboard = app.makeKeyboard(data, 2)
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
func (app *App) handleBuildProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	reg := regexp.MustCompile(buildReg)
	match := reg.FindStringSubmatch(args)
	if match == nil {
		return
	}
	_, err := app.jenkins.BuildJob(match[1])
	if err != nil {
		fmt.Printf("build error: %s", err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("build %s started", match[1]))
	msg.ParseMode = tgbotapi.ModeMarkdown
	app.bot.Send(msg)
}

const buildReg = "build ([\\w\\.\\-_\\/ ]+)?"

/**
 * handle project
 */
func (app *App) handleProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	reg := regexp.MustCompile(buildReg)
	match := reg.FindStringSubmatch(args)
	isBuild := match != nil
	switch true {
	case isBuild:
		app.handleBuildProject(message)
	case app.operators.List == args:
		app.handleListProject(message)
	case "" == args:
		app.handleListProjectStatus(message)
	}
}
