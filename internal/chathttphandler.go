package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// HttpClientManager this extends the functionality of ChatManager but
type HttpClientManager struct {
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

type httpRequestBody struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	RoomID string `json:"room_id"`
}

func (a *HttpClientManager) Message(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "not able to read the request\n")
	}
	httpRequ := httpRequestBody{}
	err = json.Unmarshal(reqBody, &httpRequ)
	if err != nil {
		fmt.Fprintf(w, "not able to parse the request")
	}
	con := &httpConnection{msg: httpRequ.Message, msgStore: a.MsgStore,userId:httpRequ.UserID,roomId:httpRequ.RoomID}
	a.ChatManager.CreateUserSession(httpRequ.UserID, con)
	w.Write([]byte("message submitted for processing\n"))

}

func (a *HttpClientManager) Command(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "not able to read the request\n")
		return
	}
	httpRequ := httpRequestBody{}
	err = json.Unmarshal(reqBody, &httpRequ)
	if err != nil {
		fmt.Fprintf(w, "not able to parse the request\n")
		return
	}
	con := &httpConnection{msg: httpRequ.Message,roomId:httpRequ.RoomID, msgStore: a.MsgStore, IsCommand: true, ch: make(chan string),userId:httpRequ.UserID}
	a.ChatManager.CreateUserSession(httpRequ.UserID, con)
	resp := <-con.ch
	w.Write([]byte(resp+"\n"))

}

func (a *HttpClientManager) Fetch(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "not able to parse the request")
		return
	}
	httpRequ := httpRequestBody{}
	err = json.Unmarshal(reqBody, &httpRequ)
	if err != nil {
		fmt.Fprintf(w, "not able to parse the request")
		return
	}
	if val, ok := a.MsgStore.store[httpRequ.UserID]; ok {
		bytePayload, _ := json.Marshal(val)
		bytePayload = append(bytePayload,[]byte("\n")...)
		w.Write(bytePayload)
		return
	}

}
