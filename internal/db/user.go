package db

type TUser struct {
	Base
	Icon int    `json:"icon,omitempty"`
	Name string `json:"name,omitempty"`
}
