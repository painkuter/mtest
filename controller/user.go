package controller

import (
	"mtest/common/errors"

	"fmt"
	"strconv"

	"time"

	"github.com/martini-contrib/sessionauth"
	"github.com/painkuter/sq"
)

const (
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

func (u UserSignUp) SaveUser() error {
	fmt.Println("Saving user")
	err := DB.Insert(&UserAuth{
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

	qb := squirrel.Update("user").
		SetMap(squirrel.Eq{"LastAccess": time.Now().String()}).
		Where(squirrel.Eq{"id_user": u.ID})

		//	str, arg,_ := qb.ToSql()
	_, err := DB.Execute(qb)

	//_, err := DB.Exec(qb.Exec())
	fmt.Println(err)
	//DB.Update(&UserAuth{
	//	ID:         u.ID,
	//	LastAccess: time.Now().String(), //TODO: fix it!
	//})
	//u.Name = name
	//u.ID = id
	fmt.Println("Logged in user " + strconv.Itoa(u.ID))
	u.authenticated = true
}

func (u *UserAuth) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserAuth) GetById(id interface{}) error {
	return nil
}

func (u *UserAuth) IsAuthenticated() bool {
	return u.authenticated
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
