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

	admin := User{LoginName: "springRain", Password: "123456", IsAdmin: "T"}

	users = append(users, admin)

}

func isAdmin(cuser User) bool {
	if cuser.IsAdmin == "T" {
		return true
	}
	return false
}
