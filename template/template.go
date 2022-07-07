package template

type Emailtemp struct {
	SENDER    string
	RECIPIENT string
	SUBJECT   string
	HTMLBODY  string
	TEXTBODY  string
	CHARSET   string
}

type Config struct {
	SLACKEN     bool
	EMAILEN     bool
	SSLEMAIL    string
	NOTEMAIL    string
	PORT        string
	SLACKHOOK   string
	CMDSTORUN   []string
	MODE        string
	MAXLOGSTORE int
}

type User struct {
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	PASSWORD []byte `json:"pwd"`
	// CREATED_AT time.Time `json:"created_at"`
}

type LoggedUser struct {
	NAME        string `json:"name"`
	EMAIL       string `json:"email"`
	ACCESSTOKEN string `json:"access-token"`
}

type SlackSchemaOne struct {
	Type     string           `json:"type"`
	Elements []SlackSchemaTwo `json:"elements"`
}

type SlackSchemaTwo struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackSchemaThree struct {
	Blocks []SlackSchemaOne `json:"blocks"`
}
