package internal

import (
	"net"
	"strings"
)

// commands which are accepted by chat application
var commands = map[string]string{
	"\\exit":                  "\t\t\tgo offline",
	"\\unsubscribe <user_id>": "block the messages from user",
	"\\rooms":                 "\t\t\tlist all the rooms/channels",
	"\\create <room_name>":    "\tcreate a new room",
	"\\join <room_name>":      "\tjoin a room",
	"\\exitroom" : "\t\tremoves users from the room",
	"\\help":      "\t\t\tprints all commands",
}

// ChatManager this is a singleton struct which spawns a go routines
// to create kind of session for chat messages to receieve and deliver
type ChatManager struct {
}

// NewChatManager returns singleton chatmanager
func NewChatManager() *ChatManager {
	return &ChatManager{}
}

// CreateUserSession creates the user session and also updates the data in chatutil struct
func (a *ChatManager) CreateUserSession(userID string, conn Conn) {
	userSession := NewUserSession(userID, conn)
	if v ,ok := conn.(*httpConnection);ok{
		userSession.room = v.roomId
	}
	chatUtil.GetSessions().AddUserSession(userSession)
	go userSession.Send()
	go userSession.Receive()
}

// CreateClientConnection starts the user first interaction
// Asks the user to enter username
// Initiase the session to send/receive the msgs
// Displays list of commands the user can use
func (a *ChatManager) CreateClientConnection(conn net.Conn) {
	conn.Write([]byte("Please enter your username ::"))
	readByte := make([]byte, 1024)
	i, err := conn.Read(readByte)
	if err != nil {
		conn.Close()
		return
	}
	con := &TCPConnection{con: conn}
	userID := strings.TrimRight(string(readByte[:i]), "\r\n")
	a.CreateUserSession(userID, con)
	conn.Write([]byte("some commands to use...\n"))
	writeHelpCommands(&TCPConnection{conn})
	conn.Write([]byte("start chatting now ...\n"))
}
