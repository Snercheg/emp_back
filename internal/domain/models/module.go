package models

type module struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ApiKey    string `json:"apikey"`
	SettingId int    `json:"setting_id"`
	TypeId    int    `json:"type_id"`
	Status    string `json:"status"`
}
