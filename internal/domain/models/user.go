package models

type User struct {
	Id       int64
	Username string
	Email    string
	PassHash []byte
}
