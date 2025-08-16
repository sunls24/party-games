package nosql

import "time"

type IMeta interface {
	GetMeta() *Meta
}

type Meta struct {
	Id      string `json:"id"`
	Version int    `json:"version,omitempty"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	VisitAt   time.Time `json:"-"`
}

func (m *Meta) GetMeta() *Meta {
	return m
}

func (m *Meta) Create(now time.Time) {
	m.CreatedAt = now
	m.Update(now)
}

func (m *Meta) Update(now time.Time) {
	m.Version++
	m.UpdatedAt = now
	m.Visit(now)
}

func (m *Meta) Visit(t time.Time) {
	m.VisitAt = t
}
