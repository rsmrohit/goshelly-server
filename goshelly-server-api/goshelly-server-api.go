package goshellyserverapi

import (
	"fmt"
	b "goshelly-server/basic"
	t "goshelly-server/template"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var r *gin.Engine

const SECRETKEY = "THIS IS A SECRET KEY, CHANGE TO SOMETHING MORE SECURE." //change this

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

func validateMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
func AddUser() {
	r.POST("/users/add/", func(c *gin.Context) {
		var user t.User
		c.BindJSON(&user)
		if !validateMailAddress(user.EMAIL) {
			c.JSON(http.StatusForbidden, gin.H{"code": "INVALID_EMAIL", "message": "Email address provided is incorrect."})
			return
		}
		if findUser(user.EMAIL) {
			c.JSON(http.StatusForbidden, gin.H{"code": "ALREADY_EXISTS", "message": "User already exists with this username."})
			return
		}
		user.PASSWORD, _ = bcrypt.GenerateFromPassword([]byte(user.PASSWORD), 12)
		os.MkdirAll("./clients/"+user.EMAIL+"/", os.ModePerm)
		f, err := os.Create("./clients/" + user.EMAIL + "/" + "pwd.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user. Service unavailable."})
			return
		}
		_, err = f.WriteString(string(user.PASSWORD))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user. Service unavailable."})
			return
		}
		f.Close()

		f, err = os.Create("./clients/" + user.EMAIL + "/" + "created_at.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User created. No creation data available."})
			return
		}
		_, err = f.WriteString(time.Now().String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User created. No creation data available."})
			return
		}
		f.Close()
		c.JSON(http.StatusCreated, gin.H{"message": "User created."})
	})
}

func findUser(uname string) bool {
	files, _ := ioutil.ReadDir("./clients")
	fmt.Println(files)
	for _, el := range files {
		if el.Name() == uname {
			return true
		}
	}
	return false
}

func RemoveUser() {
	r.DELETE("/users/remove/", func(c *gin.Context) {

		var user t.LoggedUser
		c.BindJSON(&user)
		if !findUser(user.EMAIL) {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "User not found."})
			return
		}
		if !authToken(user) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Permission denied. Ensure you are logged in."})
			return
		}
		os.RemoveAll("./clients/" + user.EMAIL)
		c.JSON(http.StatusOK, gin.H{"message": "User Deleted."})
	})
}

func LoginUser() {
	r.GET("/users/login/", func(c *gin.Context) {
		var user t.User
		c.BindJSON(&user)
		if !findUser(user.EMAIL) {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Incorrect credentials or user does not exist."})
			return
		}
		b, err := os.ReadFile("./clients/" + user.EMAIL + "/" + "pwd.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not login user. Service unavailable."})
			return
		}
		if err := bcrypt.CompareHashAndPassword(b, user.PASSWORD); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Credentials",
			})
			return
		}
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    "GoShelly Admin",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Audience:  user.EMAIL,
			Subject:   user.NAME,
		})
		token, err := claims.SignedString([]byte(SECRETKEY))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Service unavailable. Could not login.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":      "Login Successful.",
			"Access token": token,
		})
	})
}

func authToken(user t.LoggedUser) bool {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(user.ACCESSTOKEN, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if err != nil || !claims.VerifyAudience(user.EMAIL, true) ||
		!claims.VerifyIssuer("GoShelly Admin", true) || !claims.VerifyExpiresAt(time.Now().Unix(), true) || claims["sub"].(string) != user.NAME{
		return false
	}
	return true
}

func CreateLink() {
	r.GET("/users/results/", func(c *gin.Context) {
		var user t.LoggedUser
		c.BindJSON(&user)
		if !findUser(user.EMAIL) {
			c.JSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Incorrect credentials or user does not exist."})
			return
		}
		if !authToken(user) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Permission denied. Please log in again."})
			return
		}



		
	})
}


func Begin(PORT string) {
	initServerApi()
	Test()
	RemoveUser()
	AddUser()
	LoginUser()
	CreateLink()
	r.Run(":" + PORT)
}
