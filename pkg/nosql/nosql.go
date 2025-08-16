package nosql

import (
	"fmt"
	"party-games/pkg/pubsub"
	"reflect"
	"time"

	"github.com/sunls24/gox"
	bolt "go.etcd.io/bbolt"
)

type nosql struct {
	db *bolt.DB
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewDB(path string) (nosql, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second * 5})
	if err != nil {
		return nosql{}, err
	}
	return nosql{db: db}, nil
}

type Table[T IMeta] struct {
	name []byte
	db   *bolt.DB
}

func NewTable[T IMeta](n nosql) (t Table[T], err error) {
	var v T
	name := []byte(reflect.TypeOf(v).Elem().Name())
	return Table[T]{name: name, db: n.db}, n.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(name)
		return err
	})
}

func (t Table[T]) Get(id string) (v T, err error) {
	err = t.db.View(func(tx *bolt.Tx) error {
		vb := tx.Bucket(t.name).Get(gox.Str2Bytes(id))
		if vb == nil {
			return nil
		}
		v, err = gox.GobUnmarshal[T](vb)
		return err
	})
	return
}

func (t Table[T]) Delete(id string) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(t.name).Delete(gox.Str2Bytes(id))
	})
}

func (t Table[T]) Visit(id string) (v T, err error) {
	var key = gox.Str2Bytes(id)
	err = t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(t.name)
		vb := bucket.Get(key)
		if vb == nil {
			return nil
		}
		v, err = gox.GobUnmarshal[T](vb)
		if err != nil {
			return err
		}

		v.GetMeta().Visit(time.Now())
		return bucket.Put(key, gox.GobMarshal[T](v))
	})
	return
}

func (t Table[T]) Create(id string, v T) error {
	if id == "" {
		return fmt.Errorf("create %s id is empty", t.name)
	}
	v.GetMeta().Id = id
	v.GetMeta().Create(time.Now())
	return t.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(t.name).Put(gox.Str2Bytes(id), gox.GobMarshal(v))
	})
}

func (t Table[T]) Update(id string, update func(T) error) error {
	key := gox.Str2Bytes(id)
	err := t.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(t.name)
		vb := bucket.Get(key)
		if vb == nil {
			return fmt.Errorf("update %s: %s not found", t.name, id)
		}
		v, err := gox.GobUnmarshal[T](vb)
		if err != nil {
			return err
		}
		err = update(v)
		if err != nil {
			return err
		}
		v.GetMeta().Update(time.Now())
		return bucket.Put(key, gox.GobMarshal(v))
	})
	if err != nil {
		return err
	}
	return nil
}

func (t Table[T]) UpdatePublish(id string, update func(T) error) error {
	var data T
	err := t.Update(id, func(v T) error {
		data = v
		return update(v)
	})
	if err != nil {
		return err
	}
	t.Publish(id, data)
	return nil
}

func (t Table[T]) All() ([]T, error) {
	var list []T
	return list, t.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(t.name).ForEach(func(k, vb []byte) error {
			v, err := gox.GobUnmarshal[T](vb)
			if err != nil {
				return err
			}
			list = append(list, v)
			return nil
		})
	})
}

func (t Table[T]) In(ids []string, fn func(id string, v T)) error {
	if len(ids) == 0 {
		return nil
	}
	return t.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(t.name)
		for _, id := range ids {
			vb := bucket.Get(gox.Str2Bytes(id))
			if vb == nil {
				continue
			}
			v, err := gox.GobUnmarshal[T](vb)
			if err != nil {
				return err
			}
			fn(id, v)
		}
		return nil
	})
}

func (t Table[T]) InMap(ids []string) (map[string]T, error) {
	var m = make(map[string]T, len(ids))
	return m, t.In(ids, func(id string, v T) {
		m[id] = v
	})
}

func (t Table[T]) subkey(id string) string {
	return gox.Bytes2Str(t.name) + id
}

func (t Table[T]) Subscribe(id string) (<-chan T, func()) {
	ch, cancel := pubsub.Subscribe(t.subkey(id))
	to := make(chan T)
	go func() {
		defer close(to)
		select {
		case data, ok := <-ch:
			if !ok {
				return
			}
			if data == nil {
				v, err := t.Get(id)
				if err != nil { // ignore
					return
				}
				to <- v
				return
			}
			v, ok := data.(T)
			if !ok {
				return
			}
			to <- v
		}
	}()
	return to, cancel
}

func (t Table[T]) Publish(id string, data any) {
	pubsub.Publish(t.subkey(id), data)
}
