package models

type User struct {
	ID       int64
	Username string
	Email    string
	PassHash string
	Status   string
	IsAdmin  bool
}
