package controller

import (
	"mtest/common/errors"

	"fmt"
	"strconv"

	"github.com/go-gorp/gorp"
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
	ID        int    `form:"id" db:"id_user,primarykey,autoincrement"`
	UserLogin string `form:"login" db:",size:50"`
	Email     string `form:"email" db:",size:100"`
	Name      string `form:"name" db:",size:300"`
	PassHash  string `form:"pass" db:",size:300"`
	//CreatedAt     time.Time `form:"-" db:""`
	//UpdatedAt     time.Time `form:"-" db:""`
	LastAccess    string `form:"last_access"`
	authenticated bool   `form:"-" db:"-"`
}

type UserSignUp struct {
	Login    string `form:"login"`
	PassHash string `form:"pass"`
}

func (u UserSignUp) SaveUser(db *gorp.DbMap) error {
	fmt.Println("Saving user")
	err := db.Insert(&UserAuth{
		UserLogin: u.Login,
		PassHash:  u.PassHash,
		//	CreatedAt: time.Now(),
		//	UpdatedAt: time.Now(),
	})
	return err
}

func GenerateAnonymousUser() sessionauth.User {
	return &UserAuth{}
}

func (u *UserAuth) Login() {
	fmt.Println("LOGIN")
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.LastAccess = last_access
	//u.Name = name
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
	fmt.Println("CHECK AUTH")
	user := UserAuth{}
	// Move to get by pk
	err := DB.SelectOne(&user, "SELECT * FROM user WHERE UserLogin = ?", u.UserLogin)
	if err != nil {
		fmt.Println(err)
	}
	if u.PassHash == user.PassHash {
		*u = user
		return "", nil
	}
	return "", errors.New("Wrong pass/login")
}
