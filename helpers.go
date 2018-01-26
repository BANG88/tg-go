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
func (app *App) makeKeyboard(data []string) *tgbotapi.ReplyKeyboardMarkup {

	var total = len(data)
	var rows = make([][]tgbotapi.KeyboardButton, total)
	for j := 0; j < total; j++ {
		rows[j] = append(rows[j], tgbotapi.NewKeyboardButton(data[j]))
	}
	var jobKeyboard = tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        rows,
		OneTimeKeyboard: true,
		Selective:       false,
		ResizeKeyboard:  true,
	}
	return &jobKeyboard
}
