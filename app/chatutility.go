package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
	"time"
)

// once is needed so as to initialize the object only once
var once sync.Once

// chatUtil is the heart of appication which takes care of loggin,storing session data
// this also stores room data such as name of rooms which user are currently using which room
var chatUtil *ChatUtility

// MessageInfo : message template that user input/output uses
type MessageInfo struct {
	From                 string
	To                   string
	Message              string
	IsRoomMessage        bool
	RoomID               string
	Timestamp            time.Time
	IsOnlyMsgFieldToSend bool   `json:"-"`
	Err                  string `json:"err,omitempty"`
}

// ChatUtility this is kind of in memory
type ChatUtility struct {
	sessionData Session
	logger      MessageLogger
	room        *Roomstorage
}

// MessageLogger interface to log the chat messages
type MessageLogger interface {
	Log(msg MessageInfo)
}

// FileLogger implementation of Messagelogger, we are loggin the messages in files
type FileLogger struct {
	logger *log.Logger
}

// NewFileLogger takes in file descriptor and creates a new logger
func NewFileLogger(f *os.File) *FileLogger {
	logger := log.New()
	logger.SetOutput(f)
	return &FileLogger{logger}
}

// GetMessageLogger returns the logger
func (c *ChatUtility) GetMessageLogger() MessageLogger {
	return c.logger
}
func (a *FileLogger) Log(msg MessageInfo) {
	a.logger.Println(fmt.Sprintf("time : %s ,isRoomMessage : %t, senderID : %s, message : %q ", msg.Timestamp.Format("2006-01-02T15:04:05"), msg.IsRoomMessage, msg.From, msg.Message))

}

// NewChatUtility create a singleton object of chat utility
func NewChatUtility(log MessageLogger) *ChatUtility {
	once.Do(func() {
		chatUtil = &ChatUtility{sessionData: NewInMemorySession(), logger: log, room: &Roomstorage{rooms: make(map[string][]string)}}
	})
	return chatUtil
}

// GetChatUtility
func GetChatUtility() *ChatUtility {
	return chatUtil
}

// Roomstorage holds room to list of users entry
type Roomstorage struct {
	rooms map[string][]string
}

// GetRoomStorage get rooms occu
func (c *ChatUtility) GetRoomStorage() *Roomstorage {
	return c.room
}
func (c *Roomstorage) GetRooms() map[string][]string {
	return c.rooms
}
func (c *Roomstorage) GetRoomNames() []string {
	rooms := make([]string, 0)
	for k, _ := range c.rooms {
		rooms = append(rooms, k)
	}
	return rooms
}
func (c *Roomstorage) AddRoom(roomName string) {
	c.rooms[roomName] = make([]string, 0)
}
func (c *Roomstorage) DeleteUser(room, userID string) {
	if val, ok := c.rooms[room]; ok {
		for i, user := range val {
			if user == userID {
				val = append(val[:i], val[i+1:]...)
				c.rooms[room] = val
				break
			}

		}
	}
}
func (c *Roomstorage) AddUser(room, userID string) {
	c.rooms[room] = append(c.rooms[room], userID)
}

type Session interface {
	AddUserSession(session *UserSession)
	DelUserSession(string)
	GetSessionData() map[string]*UserSession
}

type InMemorySession struct {
	sessionDB map[string]*UserSession
}

func (a *ChatUtility) GetSessions() Session {
	return a.sessionData
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{sessionDB: make(map[string]*UserSession)}
}
func (a *InMemorySession) AddUserSession(session *UserSession) {
	a.sessionDB[session.userId] = session
}

func (a *InMemorySession) DelUserSession(userID string) {
	delete(a.sessionDB, userID)
}

func (a *InMemorySession) GetSessionData() map[string]*UserSession {
	return a.sessionDB
}


func formatRooms(rooms []string) string {
	var s bytes.Buffer
	for i, val := range rooms {
		s.WriteString(strconv.Itoa(i+1) + ". " + val + "\n")
	}
	return s.String()
}

// getUserToSendMsg this gets users the logic is to determine if the message is to be sent to the room users
// or all the users which are not connected to any rooms
func getUserToSendMsg(roomID string) []*UserSession {
	userSession := make([]*UserSession, 0)
	if len(roomID) != 0 {
		for _, val := range GetChatUtility().room.rooms[roomID] {
			userSession = append(userSession, chatUtil.GetSessions().GetSessionData()[val])
		}
		return userSession
	}
	for _, val := range chatUtil.GetSessions().GetSessionData() {
		if len(val.room) == 0 {
			userSession = append(userSession, val)
		}
	}
	return userSession
}

func writeHelpCommands(conn Conn) {
	msgInfo := MessageInfo{IsOnlyMsgFieldToSend: true}
	if _, ok := conn.(*httpConnection); ok {
		payload, _ := json.Marshal(commands)
		msgInfo.Message = string(payload)
		conn.Write(msgInfo)
		return
	}
	for k, v := range commands {
		conn.Write(MessageInfo{IsOnlyMsgFieldToSend: true, Message: k + ":" + v})
	}
}
