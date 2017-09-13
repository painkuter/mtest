package controller

import (
	"mtest/common/errors"

	"fmt"
	"strconv"

	"github.com/martini-contrib/sessionauth"
)

const (
	admin       = "admin"
	pass        = "1234"
	name        = "Istribitelko"
	id          = 272
	last_access = "yesterday"
)

type UserAuth struct {
	ID            int    `form:"id" db:"id_user,primarykey,autoincrement"`
	UserLogin     string `form:"login" db:",size:50"`
	Name          string `form:"name" db:",size:300"`
	PassHash      string `form:"pass" db:",size:300"`
	LastAccess    string `form:"last_access"`
	authenticated bool   `form:"-" db:"-"`
}

/*type User struct {
	ID            int    `form:"id" db:"id_user"`
	UserLogin     string `form:"login" db:"user_login"`
	Name          string `form:"name" db:"user_name"`
	PassHash      string `form:"pass" db:"user_pass"`
	LastAccess    string `form:"last_access" db:"user_last_access"`
}*/

type UserSignUp struct {
	Login    string `form:"login"`
	PassHash string `form:"pass"`
}

func (u UserSignUp) SaveUser() {

}

func GenerateAnonymousUser() sessionauth.User {
	return &UserAuth{}
}

func (u *UserAuth) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.LastAccess = last_access
	u.Name = name
	u.ID = id
	fmt.Println("Logged in user " + strconv.Itoa(u.ID))
	u.authenticated = true
}

func (u *UserAuth) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserAuth) GetById(id interface{}) error {
	//if id != 1 {
	//	return errors.New("No ID")
	//}
	return nil
}

func (u *UserAuth) IsAuthenticated() bool {
	fmt.Print("USER: ")
	fmt.Println(u)
	return u.authenticated
	//return true
}

func (u *UserAuth) UniqueId() interface{} {
	return u.ID
}

func (u *UserAuth) CheckAuth() (string, error) {

	if u.UserLogin == admin && u.PassHash == pass {
		return "ADMINKO", nil
	}
	return "", errors.New("Wrong pass/login")
}
