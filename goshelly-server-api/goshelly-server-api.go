package goshellyserverapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	b "goshelly-server/basic"
	t "goshelly-server/template"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var r *gin.Engine
var IP = getIP()
var DOMAIN = "http://"+IP+":9000"
const SECRETKEY = `U2hyaSBHdXJ1IENoYXJhbiBTYXJvb2phLXJhak5pamEgbWFudSBNdWt1cmEgU3VkaGFhcmlCYXJhbmF1IFJhaHViaGFyYSBCaW1hbGEgWWFzaHVKbyBEYXlha2EgUGhhbGEgQ2hhcmlCdWRoZWUtSGVlbiBUaGFudSBKYW5uaWtheQpTdW1pcm93IFBhdmFuYSBLdW1hcmFCYWxhLUJ1ZGhlZSBWaWR5YSBEZWhvbyBNb2hlZUhhcmFodSBLYWxlc2hhIFZpa2FhcmEuLi5KYWkgSGFudW1hbiBneWFuIGd1biBzYWdhckphaSBLYXBpcyB0aWh1biBsb2sgdWphZ2FyUmFtIGRvb3QgYXR1bGl0IGJhbCBkaGFtYUFuamFhbmktcHV0cmEgUGF2YW4gc3V0IG5hbWEuLi4KTWFoYWJpciBCaWtyYW0gQmFqcmFuZwpLdW1hdGkgbml2YXIgc3VtYXRpIEtlIHNhbmdpCkthbmNoYW4gdmFyYW4gdmlyYWogc3ViZXNhCkthbmFuIEt1bmRhbCBLdW5jaGl0IEtlc2gKSGF0aCBWYWpyYSBBdXIgRGh1dmFqZSBWaXJhagpLYWFuZGhlIG1vb25qIGphbmVodSBzYWphaQpTYW5rYXIgc3V2YW4ga2VzcmkgTmFuZGFuClRlaiBwcmF0YWFwIG1haGEgamFnIHZhbmRhbi4uLgpWaWR5YXZhYW4gZ3VuaSBhdGkgY2hhdHVyClJhbSBrYWoga2FyaWJlIGtvIGFhdHVyClByYWJ1IGNoYXJpdHJhIHN1bmliZS1rbyByYXNpeWEKUmFtIExha2huIFNpdGEgbWFuIEJhc2l5YQpTdWtzaG1hIHJvb3AgZGhhcmkgU2l5YWhpIGRpa2hhdmEKVmlrYXQgcm9vcCBkaGFyaSBsYW5rIGphcmF2YQpCaGltYSByb29wIGRoYXJpIGFzdXIgc2FuZ2hhcmUKUmFtYWNoYW5kcmEga2Uga2FqIHNhbnZhcmUuLi4KTGF5ZSBTYW5qaXZhbiBMYWtoYW4gSml5YXllClNocmkgUmFnaHV2aXIgSGFyYXNoaSB1ciBsYXllClJhZ2h1cGF0aSBLaW5oaSBiYWh1dCBiYWRhaQpUdW0gbWFtIHByaXllIEJoYXJhdC1oaS1zYW0gYmhhaQpTYWhhcyBiYWRhbiB0dW1oYXJvIHlhc2ggZ2FhdmUKQXNhLWthaGkgU2hyaXBhdGkga2FudGggbGFnYWF2ZQpTYW5rYWRoaWsgQnJhaG1hYWRpIE11bmVlc2EKTmFyYWQtU2FyYWQgc2FoaXQgQWhlZXNhLi4uCllhbSBLdWJlciBEaWdwYWFsIEphaGFuIHRlCkthdiBrb3ZpZCBrYWhpIHNha2Uga2FoYW4gdGUKVHVtIHVwa2FyIFN1Z3JlZXZhaGluIGtlZW5oYQpSYW0gbWlsYXllIHJhanBhZCBkZWVuaGEKVHVtaGFybyBtYW50cmEgVmliaGVlc2hhbiBtYWFuYQpMYW5rZXNoYXIgQmhheWUgU3ViIGphZyBqYW5hCll1ZyBzYWhhc3RyYSBqb2phbiBwYXIgQmhhbnUKTGVlbHlvIHRhaGkgbWFkaHVyIHBoYWwgamFudS4uLgpQcmFiaHUgbXVkcmlrYSBtZWxpIG11a2ggbWFoZWUKCkphbGFkaGkgbGFuZ2hpIGdheWUgYWNocmFqIG5haGVlCgpEdXJnYWFtIGthaiBqYWdhdGgga2UgamV0ZQoKU3VnYW0gYW51Z3JhaGEgdHVtaHJlIHRldGUKClJhbSBkd2FhcmUgdHVtIHJha2h2YXJlCgpIb2F0IG5hIGFneWEgYmludSBwYWlzYXJlCgpTdWIgc3VraCBsYWhhZSB0dW1oYXJpIHNhciBuYQoKVHVtIHJha3NoYWsga2FodSBrbyBkYXIgbmFhLi4uCgpBYXBhbiB0ZWogc2FtaGFybyBhYXBhaQoKVGVlbmhvbiBsb2sgaGFuayB0ZSBrYW5wYWkKCkJob290IHBpc2FhY2ggTmlrYXQgbmFoaW4gYWF2YWkKCk1haGF2aXIgamFiIG5hYW0gc3VuYXZhZQoKTmFzZSByb2cgaGFyYWUgc2FiIHBlZXJhCgpKYXBhdCBuaXJhbnRhciBIYW51bWFudCBiZWVyYQoKU2Fua2F0IHNlIEhhbnVtYW4gY2h1ZGF2YWUKCk1hbiBLYXJhbSBWYWNoYW4gZHlhbiBqbyBsYXZhaS4uLgoKU2FiIHBhciBSYW0gdGFwYXN2ZWUgcmFqYQoKVGluIGtlIGthaiBzYWthbCBUdW0gc2FqYQoKQXVyIG1hbm9yYXRoIGpvIGtvaSBsYXZhaQoKU29oaSBhbWl0IGplZXZhbiBwaGFsIHBhdmFpCgpDaGFyb24gWXVnIHBhcnRhcCB0dW1oYXJhCgpIYWkgcGVyc2lkaCBqYWdhdCB1aml5YXJhCgpTYWRodSBTYW50IGtlIHR1bSBSYWtod2FyZQoKQXN1ciBuaWthbmRhbiBSYW0gZHVsaGFyZS4uLgoKQXNodGEtc2lkaGkgbmF2IG5pZGhpIGtlIGRoYXRhCgpBcy12YXIgZGVlbiBKYW5raSBtYXRhCgpSYW0gcmFzYXlhbiB0dW1oYXJlIHBhc2EKClNhZGEgcmFobyBSYWdodXBhdGkga2UgZGFzYQoKVHVtaGFyZSBiaGFqYW4gUmFtIGtvIHBhdmFpCgpKYW5hbS1qYW5hbSBrZSBkdWtoIGJpc3JhYXZhaQoKQW50aC1rYWFsIFJhZ2h1dmlyIHB1ciBqYXllZQoKSmFoYW4gamFuYW0gSGFyaS1CYWtodCBLYWhheWVlLi4uCgpBdXIgRGV2dGEgQ2hpdCBuYSBkaGFyZWhpCgpIYW51bWFudGggc2UgaGkgc2FydmUgc3VraCBrYXJlaGkKClNhbmthdCBrYXRlLW1pdGUgc2FiIHBlZXJhCgpKbyBzdW1pcmFpIEhhbnVtYXQgQmFsYmVlcmEKCkphaSBKYWkgSmFpIEhhbnVtYW4gR29zYWhpbgoKS3JpcGEgS2FyYWh1IEd1cnVkZXYga2kgbnlhaGluCgpKbyBzYXQgYmFyIHBhdGgga2FyZSBrb3lpCgpDaHV0ZWhpIGJhbmRoaSBtYWhhIHN1a2ggaG95aS4uLgoKSm8geWFoIHBhZGhlIEhhbnVtYW4gQ2hhbGlzYQoKSG95ZSBzaWRkaGkgc2FraGkgR2F1cmVlc2EKClR1bHNpZGFzIHNhZGEgaGFyaSBjaGVyYQoKS2VlamFpIE5hdGggSHJpZGF5ZSBtZWluIGRlcmEuLi4KCktlZWphaSBOYXRoIEhyaWRheWUgbWVpbiBkZXJhLi4uCgpQYXZhbiBUYW5heSBTYW5rYXQgSGFyYW5hLApNYW5nYWxhIE11cmF0aSBSb29wLi4uClJhbSBMYWtoYW5hIFNpdGEgU2FoaXRhCkhyaWRheSBCYXNhaHUgU29vciBCaG9vcC4=` //change this


func getIP() string{
	return ""
} //fix
func initServerApi() {
	r = gin.Default()
	r.LoadHTMLGlob("html/*.html")
	os.MkdirAll("./clients/", os.ModePerm)
	os.MkdirAll("./logs/GoShellyServer-api-logs/", os.ModePerm)
	apifile, err := os.OpenFile("./logs/GoShellyServer-api-logs/"+"api-log"+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Api log open error: %s. Logs unavailable.", err)
	}
	if err == nil {
		gin.DefaultWriter = apifile
	}

	b.LogClean("./logs/GoShellyServer-api-logs/", b.SERVCONFIG.SERVMAXLOGSTORE)
}

func test() {
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
func addUser() {
	r.POST("/signup/", func(c *gin.Context) {
		var user t.User
		c.BindJSON(&user)
		if !validateMailAddress(user.EMAIL) {
			c.JSON(http.StatusForbidden, gin.H{"message": "Email address provided is incorrect."})
			return
		}
		if b.FindUser(strings.TrimSpace(user.EMAIL)) {
			c.JSON(http.StatusForbidden, gin.H{"message": "User already exists with this email. Try a different email."})
			return
		}
		user.PASSWORD, _ = bcrypt.GenerateFromPassword([]byte(user.PASSWORD), 12)
		os.MkdirAll("./clients/"+user.EMAIL+"/logs/", os.ModePerm)
		f, err := os.Create("./clients/" + user.EMAIL + "/" + "user.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user."})
			return
		}
		f.Close()
		file, err := json.MarshalIndent(t.User{
			NAME:     base64.StdEncoding.EncodeToString([]byte(user.NAME)),
			EMAIL:    base64.StdEncoding.EncodeToString([]byte(user.EMAIL)),
			PASSWORD: user.PASSWORD,
		}, "", " ")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user."})
			return
		}
		err = ioutil.WriteFile("./clients/"+user.EMAIL+"/user.json", file, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user."})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User created."})
	})
}

func removeUser() {
	r.DELETE("/delete/", func(c *gin.Context) {
		var user t.LoggedUser
		c.BindJSON(&user)
		if !b.FindUser(strings.TrimSpace(user.EMAIL)) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}
		if !authToken(user) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Permission denied. Ensure you are logged in."})
			return
		}
		os.Remove("./clients/" + user.EMAIL + "/")
		c.JSON(http.StatusOK, gin.H{"message": "User Deleted."})
	})
}

func loginUser() {
	r.POST("/login/", func(c *gin.Context) {
		var user t.LoginUser
		c.BindJSON(&user)
		if !b.FindUser(strings.TrimSpace(user.EMAIL)) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect credentials or user does not exist.", "token": ""})
			return
		}

		var temp t.User
		file, _ := ioutil.ReadFile("./clients/" + user.EMAIL + "/user.json")
		err := json.Unmarshal([]byte(file), &temp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not login user. Service unavailable.",
				"token": ""})
			return
		}
		if err := bcrypt.CompareHashAndPassword(temp.PASSWORD, user.PASSWORD); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Credentials.",
				"token":   "",
			})
			return
		}
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    "GoShelly Admin",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Audience:  user.EMAIL,
		})
		token, err := claims.SignedString([]byte(SECRETKEY))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Service unavailable. Could not login.",
				"token":   "",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Login Successful.",
			"token":   token,
		})
	})
}

func authToken(user t.LoggedUser) bool {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(user.TOKEN, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if err != nil || !claims.VerifyAudience(user.EMAIL, true) ||
		!claims.VerifyIssuer("GoShelly Admin", true) || !claims.VerifyExpiresAt(time.Now().Unix(), true) { //} || claims["sub"].(string) != user.NAME {
		return false
	}
	return true
}

func checkCurrentToken() {
	r.POST("/auth/", func(c *gin.Context) {
		var user t.LoggedUser
		c.BindJSON(&user)

		if !b.FindUser(strings.TrimSpace(user.EMAIL)) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect credentials or user does not exist.", "token": user.TOKEN})
			return
		}
		if !authToken(user) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Credentials.",
				"token":   user.TOKEN,
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Credentials=Valid",
			"token":   user.TOKEN,
		})
	})
}

func returnUserLogs() {
	r.POST("/list/", func(c *gin.Context) {
		var user t.LoggedUser
		c.BindJSON(&user)
		if !b.FindUser(strings.TrimSpace(user.EMAIL)) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No logs found for current logged in user."})
			return
		}
		if !authToken(user) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Credentials.",
			})
			return
		}
		var returnMsg strings.Builder

		files, err := ioutil.ReadDir("./clients/" + user.EMAIL + "/logs/")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get logs. Try again later."})
			return
		}
		returnMsg.WriteString("ID\t\t\tFILENAME\n")
		for id, file := range files {
			returnMsg.WriteString(strconv.Itoa(id+1) + "-->" + strings.ReplaceAll(file.Name(), ".log", "") + "\n")
		}
		c.JSON(http.StatusOK, gin.H{"message": returnMsg.String()})
	})

}

func createLink() {
	r.POST("/link/", func(c *gin.Context) {
		var user t.UserLinks
		c.BindJSON(&user)
		if !b.FindUser(user.EMAIL) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect credentials or user does not exist."})
			return
		}
		if !authToken(t.LoggedUser{
			TOKEN: user.TOKEN,
			EMAIL: user.EMAIL,
		}) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Permission denied. Please log in again."})
			return
		}
		link := DOMAIN + "/logs/" + user.EMAIL + "/" + strconv.Itoa(user.LOGID) + "/"
		c.JSON(http.StatusOK, gin.H{"message": link})
	})
}

func hostLog() {

	r.GET("/logs/:userid/:id", func(c *gin.Context) {
		userid := c.Param("userid")
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil || userid == "" || id < 1 || id > b.SERVCONFIG.CLIMAXLOGSTORE {
			c.HTML(http.StatusNotFound, "404.html", gin.H{
				"message": "Not found.",
			})
			return
		}
		files, err := ioutil.ReadDir("./clients/" + userid + "/logs/")
		if err != nil {
			c.HTML(http.StatusInternalServerError, "oops.html", gin.H{
				"message": "InternalServerError",
			})
			return
		}
		if len(files) == 0 {
			c.HTML(http.StatusNotFound, "404.html", gin.H{
				"message": "Not found.",
			})
			return
		}

		message, err := ioutil.ReadFile("./clients/" + userid + "/logs/" + files[id-1].Name())
		if err != nil {
			c.HTML(http.StatusInternalServerError, "oops.html", gin.H{
				"message": "InternalServerError",
			})
		}

		c.String(http.StatusOK, string(message))
	})
}

func startAPI() {
	initServerApi()
	removeUser()
	checkCurrentToken()
	addUser()
	loginUser()
	returnUserLogs()
	createLink()
	hostLog()
	test()
}

func BeginAPI(APIHOSTPORT string) {
	startAPI()
	r.Run(":" + APIHOSTPORT)
}
