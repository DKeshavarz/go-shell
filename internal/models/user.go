package models

import "errors"

type User struct {
	Username string `db:"username" `
	Password string `db:"password"`
}

func (u *User)Validate()(err error){
	if u.Username == ""{
		return errors.New("expeted 'user' not be empty")
	}
	return nil
}
