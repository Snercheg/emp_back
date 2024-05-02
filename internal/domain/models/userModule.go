package models

type UserModule struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	User     *User   `json:"user"`
	ModuleID int     `json:"module_id"`
	Module   *Module `json:"module"`
	Status   string  `json:"status"`
}
