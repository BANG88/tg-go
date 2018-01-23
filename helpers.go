package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/olekukonko/tablewriter"
)

/**
 * make table
 * return a formatted table or empty string
 */
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

/**
 * makeKeyboard make keyboard
 */
func (app *App) makeKeyboard(data []string, rowOfButtons int) *tgbotapi.ReplyKeyboardMarkup {
	if rowOfButtons == 0 {
		rowOfButtons = 2
	}
	var total = len(data) / rowOfButtons
	var rows = make([][]tgbotapi.KeyboardButton, total)
	var keyboard []tgbotapi.KeyboardButton
	for j := 0; j < len(data); j++ {
		keyboard = append(keyboard, tgbotapi.NewKeyboardButton(data[j]))
	}
	for index := 0; index < total; index++ {
		rows[index] = keyboard[index*rowOfButtons : index*rowOfButtons+rowOfButtons]
	}
	var jobKeyboard = tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        rows,
		OneTimeKeyboard: true,
		Selective:       false,
		ResizeKeyboard:  true,
	}
	return &jobKeyboard
}
