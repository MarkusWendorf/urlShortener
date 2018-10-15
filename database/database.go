package database

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

var (
	ErrShorthandDoesNotExist = errors.New("shorthand does not exist")
)

var bucketName = []byte("mapping")

type Database struct {
	boltDb *bolt.DB
}


func MustConnect(path string) Database {

	db, err := bolt.Open(path, 777, &bolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucketName)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return Database{boltDb: db}
}

func (db *Database) GetURL(shorthand string) (string, error) {

	var url []byte
	fmt.Println(db.boltDb)
	err := db.boltDb.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(bucketName)
		url = bucket.Get([]byte(shorthand))

		return nil
	})

	if url == nil {
		return "", ErrShorthandDoesNotExist
	}

	return string(url), err
}

func (db *Database) SetURL(shorthand string, url string) error {

	err := db.boltDb.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(bucketName)
		return bucket.Put([]byte(shorthand), []byte(url))
	})

	return err
}

func (db *Database) GetAll() {

	db.boltDb.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(bucketName)
		cursor := bucket.Cursor()
		if cursor == nil {
			return errors.New("can't find bucket named: " + string(bucketName))
		}

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			fmt.Println(string(k), string(v))
		}

		return nil
	})
}
