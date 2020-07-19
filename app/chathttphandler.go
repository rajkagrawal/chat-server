package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// HTTPClientManager this extends the functionality of ChatManager but
type HTTPClientManager struct {
	*ChatManager
	MsgStore *WebMessageStore
}

func NewWebStore() *WebMessageStore {
	return &WebMessageStore{lock: &sync.RWMutex{}, store: make(map[string][]MessageInfo)}
}

// WebMessageStore kind of stores session data for the user until the user exits the session
type WebMessageStore struct {
	lock  *sync.RWMutex
	store map[string][]MessageInfo
}

// chatReqBody  this is the common struct used to unmarshal users message
type chatReqBody struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	RoomID  string `json:"room_id"`
}

// Message is to post the chat message
// userID Message are must, if user has joined the room then roomID is must
func (a *HTTPClientManager) Message(w http.ResponseWriter, r *http.Request) {
	clientReq, err := getReqData(r, w)
	if err != nil {
		return
	}
	con := &httpConnection{msg: clientReq.Message, msgStore: a.MsgStore, userID: clientReq.UserID, roomID: clientReq.RoomID}
	a.ChatManager.CreateUserSession(clientReq.UserID, con)
	w.Write([]byte("message submitted \n"))

}

// Command when user wants to run command this api gets hit
func (a *HTTPClientManager) Command(w http.ResponseWriter, r *http.Request) {
	clientReq, err := getReqData(r, w)
	if err != nil {
		return
	}
	con := &httpConnection{msg: clientReq.Message, roomID: clientReq.RoomID, msgStore: a.MsgStore, IsCommand: true, ch: make(chan string), userID: clientReq.UserID}
	a.ChatManager.CreateUserSession(clientReq.UserID, con)
	resp := <-con.ch
	w.Write([]byte(resp + "\n"))

}

// Fetch fetches the chat message for an user, kind of user sessions
func (a *HTTPClientManager) Fetch(w http.ResponseWriter, r *http.Request) {
	clientReq, err := getReqData(r, w)
	if err != nil {
		return
	}
	if val, ok := a.MsgStore.store[clientReq.UserID]; ok {
		bytePayload, _ := json.Marshal(val)
		bytePayload = append(bytePayload, []byte("\n")...)
		w.Write(bytePayload)
		return
	}

}

func getReqData(r *http.Request, w http.ResponseWriter) (*chatReqBody, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "not able to read the request\n")
		return nil, err
	}
	clientReq := &chatReqBody{}
	err = json.Unmarshal(reqBody, &clientReq)
	if err != nil {
		fmt.Fprintf(w, "not able to parse the request")
	}
	return clientReq, nil
}
