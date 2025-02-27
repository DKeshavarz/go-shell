package models

type User struct {
	ID       int    `db:"id"`
	UserName string `db:"userame"`
	Password string `db:"password"`
}
