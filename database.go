package pictures

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/abates/disgo"
	"github.com/dgraph-io/badger"
)

var (
	ErrNotFound = errors.New("Key not found")
)

type DB interface {
	Get(bucket string, key []byte, value encoding.BinaryUnmarshaler) error
	GetValue(bucket string, key []byte) ([]byte, error)

	Put(bucket string, key []byte, value encoding.BinaryMarshaler) error
	PutValue(bucket string, key []byte, value []byte) error
}

type BadgerDB struct {
	backend *badger.DB
}

func OpenBadger(dir string) (*BadgerDB, error) {
	var err error
	db := &BadgerDB{}
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		if err == nil {
			err = fmt.Errorf("dir must be a directory not a file")
		}
		return db, err
	}

	options := badger.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	db.backend, err = badger.Open(options)
	return db, err
}

func (db *BadgerDB) Get(bucket string, key []byte, value encoding.BinaryUnmarshaler) error {
	buf, err := db.GetValue(bucket, key)
	if err == nil {
		err = value.UnmarshalBinary(buf)
	}
	return err
}

func (db *BadgerDB) GetValue(bucket string, key []byte) (buf []byte, err error) {
	key = append([]byte(bucket), key...)
	err = db.backend.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == nil {
			buf, err = item.Value()
		}
		return err
	})
	return buf, err
}

func (db *BadgerDB) Put(bucket string, key []byte, value encoding.BinaryMarshaler) error {
	buf, err := value.MarshalBinary()
	if err == nil {
		err = db.PutValue(bucket, key, buf)
	}
	return err
}

func (db *BadgerDB) PutValue(bucket string, key []byte, value []byte) error {
	key = append([]byte(bucket), key...)
	return db.backend.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (db *BadgerDB) Close() error { return db.backend.Close() }

type DisgoDB struct {
	*disgo.DB
	imgdb  DB
	ticker *time.Ticker
	done   chan chan error
	dirty  bool
}

func OpenDisgoDB(db DB) (*DisgoDB, error) {
	ddb := &DisgoDB{
		DB:     disgo.New(),
		imgdb:  db,
		ticker: time.NewTicker(time.Second * 10),
		done:   make(chan chan error),
		dirty:  false,
	}

	err := db.Get("disgo", []byte("db"), ddb)

	if err == nil || err == badger.ErrKeyNotFound {
		err = nil
		go func() {
			select {
			case ch := <-ddb.done:
				if ddb.dirty {
					ddb.imgdb.Put("disgo", []byte("db"), ddb.DB)
					ddb.dirty = false
				}
				ch <- nil
				return
			case <-ddb.ticker.C:
				if ddb.dirty {
					ddb.imgdb.Put("disgo", []byte("db"), ddb.DB)
					ddb.dirty = false
				}
			}
		}()
	}

	return ddb, err
}

func (ddb *DisgoDB) Find(hash disgo.PHash) (filename string, err error) {
	buf, err := hash.MarshalBinary()
	if err == nil {
		var value []byte
		value, err = ddb.imgdb.GetValue("disgo", buf)
		filename = string(value)
	}
	return
}

func (ddb *DisgoDB) AddHash(hash disgo.PHash, filename string) error {
	err := ddb.DB.AddHash(hash)
	if err == nil {
		ddb.dirty = true
		buf, err := hash.MarshalBinary()
		if err == nil {
			err = ddb.imgdb.PutValue("disgo", buf, []byte(filename))
		}
		if err == nil {
			err = ddb.imgdb.Put("disgo", []byte("disgo"), ddb.DB)
		}
	}
	return err
}

func (ddb *DisgoDB) Close() (err error) {
	ch := make(chan error)
	ddb.ticker.Stop()
	ddb.done <- ch
	return <-ch
}
