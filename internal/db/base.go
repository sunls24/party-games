package db

import "time"

type IBase interface {
	GetId() string
	Create()
	Update(t time.Time)
	Visit(t time.Time)
}

var _ IBase = (*Base)(nil)

type Base struct {
	Id        string     `json:"id,omitempty"`
	Version   int        `json:"version,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	VisitAt   *time.Time `json:"visitAt,omitempty"`
}

func (b *Base) GetId() string {
	return b.Id
}

func (b *Base) Create() {
	now := time.Now()
	b.CreatedAt = &now
	b.Update(now)
}

func (b *Base) Update(t time.Time) {
	b.Version++
	b.UpdatedAt = &t
	b.Visit(t)
}

func (b *Base) Visit(t time.Time) {
	b.VisitAt = &t
}

func (b *Base) R() {
	b.CreatedAt = nil
	b.UpdatedAt = nil
	b.VisitAt = nil
}
