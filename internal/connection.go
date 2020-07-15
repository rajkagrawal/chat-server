package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Conn : connection interface on which use can either exit,write or read message or data
type Conn interface {
	Close()
	Read() (bool, string, error)
	Write(info MessageInfo) error
}

// httpConnection implementation of connection interface
type httpConnection struct {
	msg       string
	userID    string
	roomID    string
	msgStore  *WebMessageStore
	resp      string
	IsCommand bool
	ch        chan string
}

func (a *httpConnection) Close() {
	delete(a.msgStore.store, a.userID)
	a.ch <- "session closed"
	return
}
func (a *httpConnection) Read() (bool, string, error) {
	return true, a.msg, nil
}
func (a *httpConnection) Write(info MessageInfo) error {
	if a.IsCommand {
		if len(info.Err) != 0 {
			a.ch <- info.Err
		} else {
			a.ch <- info.Message
		}
		a.IsCommand = false
		return nil
	}
	a.msgStore.lock.Lock()
	defer a.msgStore.lock.Unlock()
	a.msgStore.store[info.To] = append(a.msgStore.store[info.To], info)
	return nil
}

// TCPConnection implementation of conn interface
// this takes care of users connecting via telnet
type TCPConnection struct {
	con net.Conn
}

// Close : closes the underlying  connection
func (a *TCPConnection) Close() {
	a.con.Close()
}

// Read : read the users input/msgs from given medium
func (a *TCPConnection) Read() (bool, string, error) {
	s, err := bufio.NewReader(a.con).ReadString('\n')
	if err != nil {
		return true, "", err
	}
	s = strings.Trim(s, "\r\n")
	return false, s, nil
}

// Write : this writes the msgs/error msgs to the users medium or terminal
// Some of the msg that are not chat messages are also sent via IsOnlyMsgFieldToSend
func (a *TCPConnection) Write(info MessageInfo) error {
	var err error
	if info.IsOnlyMsgFieldToSend {
		_, err = a.con.Write([]byte(fmt.Sprintf("%s\n", info.Message)))
	} else {
		_, err = a.con.Write([]byte(fmt.Sprintf("time : %s, senderID : %s, msg : %s\n", info.Timestamp.Format("2006-01-02T15:04:05"), info.From, info.Message)))
	}
	return err
}
