package app

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	systemUser = "system"
)

// Below are various regex which equates commands which are sent as message by user
var joinRegex = regexp.MustCompile(`^\\join+\s[a-z]+$`)
var createRegex = regexp.MustCompile(`^\\create+\s[a-z]+$`)
var unsubscribeRegex = regexp.MustCompile(`^\\unsubscribe+\s[a-z]+$`)

// UserSession keeps user data, holds the connection obj on which to send the messages
// msg where we decouple message send
type UserSession struct {
	conn         Conn
	msg          chan MessageInfo
	userId       string
	created      time.Time
	room         string
	unsubscribed map[string]struct{}
}

// NewuserSession creates the user session object holding connection object and initializing few other parameters required for user session
func NewUserSession(userID string, conn Conn) *UserSession {
	return &UserSession{conn: conn, userId: userID, msg: make(chan MessageInfo), created: time.Now(), unsubscribed: make(map[string]struct{})}
}

// Unsubscribe this puts the user to unsubscribe in user object so that users can be skipped while sending the message
func (u *UserSession) Unsubscribe(user string) error {
	u.unsubscribed[user] = struct{}{}
	return nil
}

// Send this is the core of this application where user command/text are read from the command line
// We run most of the code in for loop to read from input so that we can capture msgs/command when the user types in
// This has few cases where command are equated against to take some particular action
func (u *UserSession) Send() {
	isExited := false
	for !isExited {
		isExit, msgToSend, err := u.ReadInput()
		if err != nil {
			panic("some error reading the messages")
		}
		if len(msgToSend) == 0 {
			isExited = isExit
			continue
		}
		msgToSend = strings.TrimSpace(msgToSend)
		msgInfo := MessageInfo{From: u.userId, Timestamp: time.Now(), Message: msgToSend, RoomID: u.room}
		switch {
		case msgToSend == "\\exit":
			u.msg <- msgInfo
			chatUtil.GetSessions().DelUserSession(u.userId)
			//check if user belonged to room
			if len(u.room) > 0 {
				chatUtil.GetRoomStorage().DeleteUser(u.room, u.userId)
			}
			msgInfo.From = systemUser
			msgInfo.Message = fmt.Sprintf("%s logged out", u.userId)
			u.conn.Close()
			isExit = true
			u.SendMessageToUsers(msgInfo)
		case unsubscribeRegex.MatchString(msgToSend):
			userToUnsubscribe := strings.Split(msgToSend, " ")[1]
			err := u.Unsubscribe(userToUnsubscribe)
			if err != nil {
				u.conn.Write(MessageInfo{Err: err.Error(), Message: err.Error()})
			} else {
				u.conn.Write(MessageInfo{Message: "user unsubscribed", IsOnlyMsgFieldToSend: true, Timestamp: time.Now()})
			}
		case msgToSend == "\\rooms":
			rooms := chatUtil.GetRoomStorage().GetRoomNames()
			formatRooms(rooms)
			u.conn.Write(MessageInfo{Timestamp: time.Now(), IsOnlyMsgFieldToSend: true, Message: formatRooms(rooms)})
		case joinRegex.MatchString(msgToSend):
			roomName := strings.Split(msgToSend, " ")[1]
			if err := u.JoinRoom(roomName); err != nil {
				u.conn.Write(MessageInfo{Err: err.Error(), Message: err.Error(), IsOnlyMsgFieldToSend: true})
			} else {

				msgInfo.From = systemUser
				msgInfo.Message = fmt.Sprintf("user : %s, joined the room ", u.userId)
				msgInfo.IsRoomMessage = true
				u.SendMessageToUsers(msgInfo)
			}
		case msgToSend == "\\help":
			writeHelpCommands(u.conn)

		case createRegex.MatchString(msgToSend):
			roomName := strings.Split(msgToSend, " ")[1]
			err := u.CreateRoom(roomName)
			if err != nil {
				u.conn.Write(MessageInfo{Err: err.Error(), IsOnlyMsgFieldToSend: true, Message: err.Error()})
			}
		case msgToSend == "\\exitroom":
			msgInfo.From = systemUser
			msgInfo.Message = fmt.Sprintf("user : %s, exited the room ", u.userId)
			msgInfo.IsRoomMessage = true
			u.ExitRoom()
			u.SendMessageToUsers(msgInfo)
			u.room = ""
			msgInfo.IsOnlyMsgFieldToSend = true
			u.conn.Write(msgInfo)

		default:
			u.SendMessageToUsers(msgInfo)
		}
		isExited = isExit
	}
}

// SendMessageToUsers this send the messages to recipeints via a channel on which users would be listening to
// this checks if the message is intended for room or not
// logs all the messages that are sent
// check if user was unsubscribed or not
func (u *UserSession) SendMessageToUsers(msgInfo MessageInfo) {
	isRommMessage := true
	if len(u.room) == 0 {
		isRommMessage = false
	}
	msgInfo.IsRoomMessage = isRommMessage
	msgInfo.RoomID = u.room
	chatUtil.GetMessageLogger().Log(msgInfo)
	for _, v := range getUserToSendMsg(u.room) {
		if _, ok := v.unsubscribed[u.userId]; !ok {
			msgInfo.To = v.userId
			v.msg <- msgInfo
		}
	}
}

// exitRoom exiting from the room so that user can come out of the room and send some individual messages
// remove the user entry from room array
func (u *UserSession) ExitRoom() {
	roomUsers := chatUtil.GetRoomStorage().GetRooms()[u.room]
	i := 0
	for ; i < len(roomUsers); i++ {
		if roomUsers[i] == u.userId {
			break
		}
	}
	roomUsers = append(roomUsers[:i], roomUsers[i+1:]...)
	chatUtil.GetRoomStorage().GetRooms()[u.room] = roomUsers

}

// JoinRoom this attaches user to the room so that user can send msgs to only the given room
// Check if the room is available or else throw error
func (u *UserSession) JoinRoom(roomName string) error {
	isRoomAvailable := false
	for _, val := range GetChatUtility().GetRoomStorage().GetRoomNames() {
		if val == roomName {
			isRoomAvailable = true
		}
	}
	if !isRoomAvailable {
		return errors.New("The room is not yet created please first create to join")
	}
	u.room = roomName
	GetChatUtility().room.AddUser(roomName, u.userId)
	return nil
}

//CreateRoom creates the room if not already exist
func (u *UserSession) CreateRoom(roomName string) error {
	for k, _ := range GetChatUtility().room.rooms {
		if k == roomName {
			return errors.New("room already exists")
		}
	}
	GetChatUtility().GetRoomStorage().AddRoom(roomName)
	u.conn.Write(MessageInfo{Timestamp: time.Now(), From: systemUser, To: u.userId, Message: fmt.Sprintf("room created %s", roomName)})
	return nil

}

// This will read the input from the given connection which user session object holds
// this has three return parameters bool parameter is to check if we have persistent connection like telnet
// or http connection which is not persistent and will no
func (u *UserSession) ReadInput() (bool, string, error) {
	return u.conn.Read()
}

// Recieve this has a channel which receives the message that an user will receive
func (u *UserSession) Receive() {
	for {
		msg := <-u.msg
		if msg.Message == "\\exit" {
			return
		}
		u.conn.Write(msg)
	}
}
