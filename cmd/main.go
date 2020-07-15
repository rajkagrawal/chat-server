package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rajaanova/chat-server/internal/config"
	"net/http"

	"github.com/rajaanova/chat-server/internal"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

func main() {
	configFile := flag.String("configfile", "./config.json", "config file location")
	flag.Parse()
	appConfig, err := config.BS{}.Boot(*configFile)
	if err != nil {
		//panicing since configuration should be properly loaded.
		panic(fmt.Sprintf("not able bootstrap the configuration %v", err))
	}
	//Create the listener to listen from telnet connection
	listener, err := net.Listen("tcp", appConfig.ServerIP+":"+appConfig.PortNum)
	if err != nil {
		panic(fmt.Sprintf("not able to listen on  %v:%v error: %v ", appConfig.ServerIP, appConfig.PortNum, err))
	}
	defer listener.Close()
	log.Info("listening on: ", listener.Addr())

	// Getting the file description of log file and setting the required flag to make it appendible and create the file if not present
	f, err := os.OpenFile(appConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Sprintf("error while adding file to logger %v", err))
	}
	// creating the objects required
	fileLogger := internal.NewFileLogger(f)
	internal.NewChatUtility(fileLogger)
	clientManager := internal.NewChatManager()
	webMsgStore := internal.NewWebStore()
	// integrating the http api to integrate with chat system
	httpHanlder := internal.HTTPClientManager{ChatManager: clientManager, MsgStore: webMsgStore}
	// various router or handlers for http api
	router := mux.NewRouter()
	router.HandleFunc("/post", httpHanlder.Message).Methods(http.MethodPost)
	router.HandleFunc("/fetch", httpHanlder.Fetch).Methods(http.MethodPost)
	router.HandleFunc("/command", httpHanlder.Command).Methods(http.MethodPost)
	// run http servier
	go func() {
		err = http.ListenAndServe(appConfig.HTTPServerIP+":"+appConfig.HTTPPortNum, router)
		if err != nil {
			panic(err)
		}
	}()
	//listen to telnet connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}
		go clientManager.CreateClientConnection(conn)
	}
}
