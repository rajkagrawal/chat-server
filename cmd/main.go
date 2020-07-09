package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"github.com/rajaanova/chat-server/internal"
)

var (
	serverIP *string
	portNum  *string
	logFile  *string
)

func init() {
	serverIP = flag.String("ip", "127.0.0.1", "ip address to listen")
	portNum = flag.String("port", "5050", "port number on which  to listen")
	logFile = flag.String("logfile", "./message.log", "log file location")
	flag.Parse()
}

const timeFormat = "02/01/2006 15:04:05"

func main() {
	fmt.Println("portNume ", *portNum)
	listener, err := net.Listen("tcp", *serverIP+":"+*portNum)
	if err != nil {
		panic(fmt.Sprintf("not able to listen on  %v:%v error: %v ", *serverIP, *portNum, err))
	}
	defer listener.Close()
	log.Println("listening on: ", listener.Addr())
	//internal.ClientManager{
	//}
	//main listen accept loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}
		go createClientConn(conn)
	}
}

func createClientConn(conn net.Conn) {
	//ask for username and password
	conn.Write([]byte("Please enter your username and password sepated by colon i.e <username>:<password>\n"))
	readByte := make([]byte, 1024)
	_, err := conn.Read(readByte)
	if err != nil {
		conn.Close()
	}


}
