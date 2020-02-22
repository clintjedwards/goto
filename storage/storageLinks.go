package storage

import (
	"bytes"
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/go/models"
)

// func (db *BoltDB) GetLink(shortURL string) (*models.Link, error) {

// 	return nil, nil
// }

// func (db *BoltDB) GetAllLinks() () {

// }
// func (db *BoltDB) AddJob(account string, newJob *api.Job) (key string, err error) {

// CreateLink stores a new link into database
func (db *BoltDB) CreateLink(link *models.Link) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(linksBucket))

		buf := &bytes.Buffer{}
		err := binary.Write(buf, binary.BigEndian, link)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(link.Name), buf.Bytes())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// func (db *BoltDB) UpdateLink() () {

// }

// func (db *BoltDB) DeleteLink() () {

// }
