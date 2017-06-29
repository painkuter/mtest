package controller

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"golang.org/x/crypto/openpgp/errors"
)

const (
	admin = "admin"
	pass  = "1234"
)

type UserAuth struct {
	ID            int    `form: "-"`
	Name          string `form: "login"`
	PassHash      string `form: "pass"`
	authenticated bool   `form: "-"`
}

func PrintUser(r render.Render) error {
	return nil
}

func GenerateAnonymousUser() sessionauth.User {
	return &UserAuth{}
}

func (u *UserAuth) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

func (u *UserAuth) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserAuth) GetById(id interface{}) error {
	if id != 1 {
		return errors.ErrUnknownIssuer
	}
	return nil
}

func (u *UserAuth) IsAuthenticated() bool {
	return u.authenticated
}

func (u *UserAuth) UniqueId() interface{} {
	return u.ID
}
