package internal

import "errors"

type ClientManager struct {
	auth Authentication

}


type  Authentication interface {
	authenticate(string,string ) (bool,error)
}

type Room struct {
	displayName string
	roomID string
	members []UserInfo
}
type UserInfo struct {
	 userName string
	 UserProfile
	 password string
}
type UserProfile struct {
	firstName string
	lastName string
	// we can add other fields but for now just adding basic fields
}
type InmemoryAuthenctication struct {
	database map[string]UserInfo
}

func (a *InmemoryAuthenctication) authenticate(userName , password string ) (bool, error) {
	if val,ok := a.database[userName];ok {
		if val.password == password{
			return true , nil
		}
	}else{
		return false, errors.New("please check your username or password or sign up to use our service ")
	}
}