# Telegram-Jenkins

[![Build Status](https://travis-ci.org/bang88/tg-go.svg?branch=master)](https://travis-ci.org/bang88/tg-go)
[![Build status](https://ci.appveyor.com/api/projects/status/1slye681x7ekaa88/branch/master?svg=true)](https://ci.appveyor.com/project/bang88/tg-go/branch/master)

> Building jenkins project from telegram app

- [Telegram-Jenkins](#telegram-jenkins)
	- [Features](#features)
	- [Requirements](#requirements)
	- [Installation](#installation)
	- [Notification](#notification)
## Features

- list all projects(/sub projects) in chat
- manage who can use this bot
- build as one click

## Requirements

* dep for install go dependencies
* bot token generate from bot father
* jenkins server

## Installation

Download one of the executable file from [Release](https://github.com/bang88/tg-go/releases)

and change permission if needed

Write some configurations:

```yaml
dbPath: bot.db
jenkins:
  server:'jenkins-server-address'
  username: jenkins_admin
  password: 'jenkins_password'
  # dont change this value. if you want get notification from jenkins server you need install a notification plugin which will use this field
  telegramChatId: 'Telegram_Chat_ID'
botToken: 'bot_token_generated_from_bot_father'
superAdmin: 'default_admin(telegram_username)'

```

## Notification

If you don't need build notification you can skip this step

- Download notification plugin from [here](https://github.com/bang88/build-notifications-plugin/releases/download/v1.5.1/build-notifications.hpi) and then install it on your jenkins server
- Notification configuration please checkout the [docs](https://github.com/bang88/build-notifications-plugin)
- ATM: you must add a parameterized build named `Telegram_Chat_ID` and leave the default value empty
- Add a post build step(Telegram Notification) in your jenkins project
- fill up the Global Notification Target as `${Telegram_Chat_ID}` received from your last step's settings as a env variable.
- done

Why need `Telegram_Chat_ID` because jenkins need to know which chat you want post message to. this variable will be replaced in tg-bot. the bot get the `chat_id` from the telegram's chat.
