package main

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// JenkinsConf Jenkins
type JenkinsConf struct {
	Server         string `yaml:"server"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	TelegramChatID string `yaml:"telegramChatId"`
}

// Conf app configurations
type Conf struct {
	DbPath     string      `yaml:"dbPath"`
	Jenkins    JenkinsConf `yaml:"jenkins"`
	BotToken   string      `yaml:"botToken"`
	SuperAdmin string      `yaml:"superAdmin"`
}

// GetConf get
func GetConf() Conf {
	var conf Conf
	const fileName = "conf.yml"
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fileName)
			defer f.Close()
			if err != nil {
				log.Fatalf("create file %s failed: %s", fileName, err)
			}
			f.WriteString(`
# bot configurations
dbPath: bot.db
jenkins:
  server: 'jenkins-server-address'
  username: jenkins_admin
  password: 'jenkins_password'
  telegramChatId: 'Telegram_Chat_ID'
botToken: 'bot_token_generated_from_bot_father'
superAdmin: 'default_admin_telegram_username'
`)
			log.Fatalf("File %s not exists. creating one for you now. \nyou should configure it by yourself then restart the bot", fileName)

		}
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return conf
}
