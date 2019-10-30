package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bndr/gojenkins"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const folder = "com.cloudbees.hudson.plugins.folder.Folder"

/**
 * handleListProject list all projects as keyboard
 */
func (app *App) handleListProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	folderName := app.getArgument(args)
	endpoint := "/"
	if folderName != "" {
		endpoint = "/job/" + folderName
	}
	var data []string
	exec := gojenkins.Executor{Raw: new(gojenkins.ExecutorResponse), Jenkins: app.jenkins}
	_, err := app.jenkins.Requester.GetJSON(endpoint, exec.Raw, nil)
	for _, job := range exec.Raw.Jobs {
		var ji *gojenkins.Job
		if folderName != "" {
			ji, _ = app.jenkins.GetJob(job.Name, folderName)
		} else {
			ji, _ = app.jenkins.GetJob(job.Name)
		}

		if err != nil {
			log.Printf("get job failure: %s", err)
			return
		}
		if ji.Raw.Class == folder {
			data = append(data, fmt.Sprintf("/%s %s %s", app.commands.Project, app.operators.List, ji.GetName()))
		} else {
			jobName := ji.GetName()
			if folderName != "" {
				jobName = fmt.Sprintf("%s/job/%s", folderName, ji.GetName())
			}
			data = append(data, fmt.Sprintf("/%s %s %s", app.commands.Project, app.operators.Build, jobName))
		}
	}
	if err != nil {
		log.Printf("get jobs failure: %s", err)
		return
	}

	var jobKeyboard = app.makeKeyboard(data)
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID
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
	var str = app.makeTable(data, nil)
	log.Printf("str: %s", str)
	msg := tgbotapi.NewMessage(message.Chat.ID, str)
	msg.ParseMode = tgbotapi.ModeMarkdown
	app.bot.Send(msg)
	app.handleListProject(message)
}

// handleBuildProject build project
func (app *App) handleBuildProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	command := app.getArgument(args)
	if command == "" {
		return
	}
	// If we want receive build results from Jenkins we
	// need add these params
	var params = map[string]string{app.conf.Jenkins.TelegramChatID: strconv.FormatInt(message.Chat.ID, 10)}
	_, err := app.jenkins.BuildJob(command, params)
	if err != nil {
		fmt.Printf("Build error: %s", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Build %s failure 😢: %s", command, err))
		msg.ParseMode = tgbotapi.ModeMarkdown
		app.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Build started for: %s 😊", command))
		msg.ParseMode = tgbotapi.ModeMarkdown
		app.bot.Send(msg)
	}

}

const buildReg = "\\w+ ([\\w\\.\\-_\\/ ]+)?"

/**
 * handle project
 */
func (app *App) handleProject(message *tgbotapi.Message) {
	args := app.getCommandArguments(message)
	isBuild := strings.HasPrefix(args, app.operators.Build)
	switch true {
	case isBuild:
		app.handleBuildProject(message)
	case "" == args:
		app.handleListProjectStatus(message)
	default:
		app.handleListProject(message)
	}
}
