package main

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		conf string
		want string
	}{
		{
			name: "dbPath",
			conf: "DbPath",
			want: "bot.db",
		},
		{
			name: "botToken",
			conf: "BotToken",
			want: "",
		},
		{
			name: "jenkins",
			conf: "Server",
			want: "http://localhost:32769/",
		},
		{
			name: "jenkins",
			conf: "Username",
			want: "bot",
		},
		{
			name: "jenkins",
			conf: "Password",
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Conf
			conf := c.getConf()
			r := reflect.ValueOf(conf)
			k := reflect.Indirect(r).FieldByName(tt.conf)
			switch k.Kind() {
			case reflect.String:
				if k.String() != tt.want {
					t.Errorf("%s: not match %s:%s", tt.conf, k.String(), tt.want)
				}
			// treat it as Jenkins field
			case reflect.Invalid:
				rj := reflect.ValueOf(conf.Jenkins)
				kj := reflect.Indirect(rj).FieldByName(tt.conf)
				if kj.String() != tt.want {
					t.Errorf("%s: not match %s:%s", tt.conf, kj.String(), tt.want)
				}

			}
		})
	}
}
