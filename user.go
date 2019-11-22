package table

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Permission int

const (
	DownloadTable Permission = iota
	EditorTable
	CreateTable
)

var (
	expTime = 3 * 24 * time.Hour
	userMux = sync.Mutex{}
	userMap = map[string]*User{}
)

type User struct {
	UserName    string
	Permissions []Permission
	Token       string
	expiredTime time.Time
}

func NewUser(userName, per string) *User {
	var per_ []Permission
	MustJsonUnmarshal(([]byte)(per), &per_)
	return &User{
		UserName:    userName,
		Permissions: per_,
		Token:       makeToken(userName),
		expiredTime: time.Now().Add(expTime),
	}
}

func makeToken(userName string) string {
	now := time.Now().Unix()
	return fmt.Sprintf("%d@%s", now, userName)
}

func processToken(token string) (string, string) {
	s := strings.Split(token, "@")
	return s[0], s[1]
}

func AddUser(user *User) {
	userMux.Lock()
	defer userMux.Unlock()
	if _, ok := userMap[user.UserName]; !ok {
		userMap[user.UserName] = user
	}
}

func RemoveUser(userName string) {
	userMux.Lock()
	defer userMux.Unlock()
	if _, ok := userMap[userName]; ok {
		delete(userMap, userName)
	}
}

func GetUser(userName string) *User {
	userMux.Lock()
	defer userMux.Unlock()
	if u, ok := userMap[userName]; ok {
		return u
	}
	return nil
}
