package db

import (
	"fmt"
	"party-games/internal/utils"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	User Table[*TUser]
	Room Table[*TRoom]
	Game Table[*TGame]
)

func Init(path string) error {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second * 5})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		User, err = NewTable[*TUser]([]byte("user"), tx)
		if err != nil {
			return err
		}
		Room, err = NewTable[*TRoom]([]byte("room"), tx)
		if err != nil {
			return err
		}
		Game, err = NewTable[*TGame]([]byte("game"), tx)
		if err != nil {
			return err
		}
		return nil
	})
}

type Table[T IBase] struct {
	name []byte
	db   *bolt.DB
}

func NewTable[T IBase](name []byte, tx *bolt.Tx) (Table[T], error) {
	_, err := tx.CreateBucketIfNotExists(name)
	return Table[T]{name: name, db: tx.DB()}, err
}

func (t Table[T]) Get(id string) (v T, err error) {
	var key = []byte(id)
	err = t.db.Update(func(tx *bolt.Tx) error {
		bu := tx.Bucket(t.name)
		b := bu.Get(key)
		if b == nil {
			return fmt.Errorf("%s %s not found", t.name, id)
		}
		v, err = utils.Unmarshal[T](b)
		if err != nil {
			return err
		}

		v.Visit(time.Now())
		return bu.Put(key, utils.Marshal[T](v))
	})
	return
}

func (t Table[T]) Create(v T) error {
	if v.GetId() == "" {
		return fmt.Errorf("create %s id is empty", t.name)
	}
	v.Create()
	return t.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(t.name).Put([]byte(v.GetId()), utils.Marshal(v))
	})
}

func (t Table[T]) Update(id string, fn func(T) error) error {
	key := []byte(id)
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(t.name).Get(key)
		if b == nil {
			return fmt.Errorf("update %s %s not found", t.name, id)
		}
		v, err := utils.Unmarshal[T](b)
		if err != nil {
			return err
		}
		err = fn(v)
		if err != nil {
			return err
		}
		v.Update(time.Now())
		return tx.Bucket(t.name).Put(key, utils.Marshal(v))
	})
}

func (t Table[T]) Exist(id string) (bool, error) {
	var v []byte
	err := t.db.View(func(tx *bolt.Tx) error {
		v = tx.Bucket(t.name).Get([]byte(id))
		return nil
	})
	return v != nil, err
}

func (t Table[T]) All() ([]T, error) {
	var list []T
	return list, t.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(t.name).ForEach(func(k, vb []byte) error {
			v, err := utils.Unmarshal[T](vb)
			if err != nil {
				return err
			}
			list = append(list, v)
			return nil
		})
	})
}
