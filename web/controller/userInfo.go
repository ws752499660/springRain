package controller

import "github.com/hqu.edu.cn/springRain/service"

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName string
	Password  string
	IsAdmin   string
}

var users []User

func init() {

	admin := User{LoginName: "Hanxiaodong", Password: "123456", IsAdmin: "T"}
	alice := User{LoginName: "springRain", Password: "123456", IsAdmin: "T"}
	bob := User{LoginName: "alice", Password: "123456", IsAdmin: "F"}
	jack := User{LoginName: "bob", Password: "123456", IsAdmin: "F"}

	users = append(users, admin)
	users = append(users, alice)
	users = append(users, bob)
	users = append(users, jack)

}

func isAdmin(cuser User) bool {
	if cuser.IsAdmin == "T" {
		return true
	}
	return false
}
