package model

type User struct {
	Userid   int64
	Username string
	Password string
	Salt     string
}
