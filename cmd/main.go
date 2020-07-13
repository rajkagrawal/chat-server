package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/rajaanova/chat-server/internal"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

// Declaring the variable that need to be passed in using command line parameters
// for e.g  go run cmd/main.go -port=8888
var (
	serverIP *string
	portNum  *string
	httpPortNum *string
	logFile  *string
)

// init : initializing the parmeters using flag utility, if any parameters is not passed then use default parameters mentioned below
func init() {
	serverIP = flag.String("ip", "127.0.0.1", "ip address to listen")
	portNum = flag.String("port", "5050", "port number on which  to listen")
	httpPortNum = flag.String("httpport", "8080", "http port number on which  to listen")
	logFile = flag.String("logfile", "./message.log", "log file location")
	flag.Parse()
}

func main() {
	//Create the listener to listen from telnet connection
	listener, err := net.Listen("tcp", *serverIP+":"+*portNum)
	if err != nil {
		panic(fmt.Sprintf("not able to listen on  %v:%v error: %v ", *serverIP, *portNum, err))
	}
	defer listener.Close()
	log.Info("listening on: ", listener.Addr())

	// Getting the file description of log file and setting the required flag to make it appendible and create the file if not present
	f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Sprintf("error while adding file to logger %v", err))
	}
	// creating the objects required
	fileLogger := internal.NewFileLogger(f)
	internal.NewChatUtility(fileLogger)
	clientManager := internal.NewChatManager()
	webMsgStore := internal.NewWebStore()
	// integrating the http api to integrate with chat system
	httpHanlder := internal.HttpClientManager{clientManager, webMsgStore}
	// various router or handlers for http api
	router := mux.NewRouter()
	router.HandleFunc("/post", httpHanlder.Message).Methods(http.MethodPost)
	router.HandleFunc("/fetch", httpHanlder.Fetch).Methods(http.MethodPost)
	router.HandleFunc("/command", httpHanlder.Command).Methods(http.MethodPost)
	// run http servier
	go func() {
		err = http.ListenAndServe(":"+*httpPortNum, router)
		if err != nil {
			panic(err)
		}
	}()
	//listen accept loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}
		go clientManager.CreateClientConnection(conn)
	}
}
