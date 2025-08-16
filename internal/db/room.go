package db

type TRoom struct {
	Base
	Owner string   `json:"owner,omitempty"`
	Users []TRUser `json:"users,omitempty"`
}

func (r *TRoom) R() *TRoom {
	r.Base.R()
	r.Id = r.Id[len(r.Id)-4:]
	for i := range r.Users {
		r.Users[i].R()
	}
	return r
}

func (r *TRoom) SetOnline(check func(string, string) bool) *TRoom {
	for i, v := range r.Users {
		if v.Id == "" {
			continue
		}
		r.Users[i].Online = check(r.Id, v.Id)
	}
	return r
}

type TRUser struct {
	TUser
	Online bool `json:"online,omitempty"`
}
