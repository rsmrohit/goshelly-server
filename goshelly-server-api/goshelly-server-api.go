package goshellyserverapi

import (
	"fmt"
	b "goshelly-server/basic"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	t "goshelly-server/template"
)



var r *gin.Engine

func initServerApi() {
	r = gin.Default()
	os.MkdirAll("./clients/", os.ModePerm)
	os.MkdirAll("./logs/GoShellyServer-api-logs/", os.ModePerm)

	apifile, err := os.OpenFile("./logs/GoShellyServer-api-logs/"+"api-log"+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Api log open error: %s. Logs unavailable.", err)
	}
	if err == nil {
		gin.DefaultWriter = apifile
	}


	b.LogClean("./logs/GoShellyServer-api-logs/", 100)
	///NOTE: 100 is a random hardcoded value, this function call decides that the max number of users cannot excede 100
	// due to memory constraints as everything is stored in memory for now.

}

func Test() {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func AddUser() {
	r.POST("/users/", func(c *gin.Context) {
		//do some authentication stuff
		var user t.User
		c.BindJSON(&user)

		if findUser(user.USERNAME){
			c.JSON(http.StatusForbidden, gin.H{"code": "ALREADY_EXISTS", "message": "User already exists with this username."})
		}
		os.MkdirAll("./clients/"+ user.USERNAME+"/",os.ModePerm)
		c.JSON(http.StatusOK,  gin.H{"message": "User created."})
	})
}
func findUser(uname string)bool {
	files, _ := ioutil.ReadDir("/clients/")
	for _, el := range files{
		if el.Name() == uname {
			return true
		}
	}
	return false
}

func RemoveUser() {
	r.DELETE("/users/remove/", func(c *gin.Context){
		
		//do some authentication stuff
		var user t.User
		c.BindJSON(&user)
		if !findUser(user.USERNAME){
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "User not found, or doesn't exist"})
		}
		os.RemoveAll("./clients/"+ user.USERNAME)
		c.JSON(http.StatusOK,  gin.H{"message": "User Deleted."})
	})
}

func AuthLoginUser() {

}

func CreateLink() {
	
}

func Begin(PORT string) {
	initServerApi()
	Test()
	RemoveUser()
	AddUser()
	r.Run(":9000")
}
