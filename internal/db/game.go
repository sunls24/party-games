package db

type TGame struct {
	Base
	Started bool `json:"started"`
}
