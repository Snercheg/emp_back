package models

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	PassHash  []byte
	Status    string
	RoleId    int64
	Role      *Role
	CreatedAt time.Duration
}
