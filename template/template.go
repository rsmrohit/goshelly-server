package template

import (
	"time"
)

type Emailtemp struct {
	SENDER    string
	RECIPIENT string
	SUBJECT   string
	HTMLBODY  string
	TEXTBODY  string
	CHARSET   string
}

type Config struct {
	SLACKEN   bool
	EMAILEN   bool
	SSLEMAIL  string
	NOTEMAIL  string
	PORT      string
	SLACKHOOK string
	CMDSTORUN []string
	MODE      string
	MAXLOGSTORE int
}


type User struct{
	USERNAME string `json:"uname"`
	PASSWORD string `json:"pwd"`
	CREATED_AT time.Time `json:"created_at"`
}

type SlackSchemaOne struct {
	Type     string           `json:"type"`
	Elements []SlackSchemaTwo `json:"elements"`
}

type SlackSchemaTwo struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackSchemaThree struct{
	Blocks []SlackSchemaOne `json:"blocks"`
}
