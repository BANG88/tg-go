package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)
// Jenkins Jenkins
type Jenkins struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Conf app configurations
type Conf struct {
	DbPath   string  `yaml:"dbPath"`
	Jenkins  Jenkins `yaml:"jenkins"`
	BotToken string  `yaml:"botToken"`
}

// GetConfig get
func (conf *Conf) getConf() *Conf {
	data, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return conf
}
