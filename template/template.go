package template



type Emailtemp struct {
	SENDER string
	RECIPIENT string
	SUBJECT string
	HTMLBODY string
	TEXTBODY string
	CHARSET string
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
}

type SlackPost struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Data  []ComRes `json:"data"`
}
type ComRes struct {
	Cmd string `json:"cmd"`
	Res string `json:"res"`
}


