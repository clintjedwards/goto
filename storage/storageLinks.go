package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/go/models"
)

// GetAllLinks returns an unpaginated list of current links
func (db *BoltDB) GetAllLinks() (map[string]models.Link, error) {

	results := map[string]models.Link{}

	db.store.View(func(tx *bolt.Tx) error {
		linksBucket := tx.Bucket([]byte(linksBucket))

		err := linksBucket.ForEach(func(key, value []byte) error {
			var link models.Link

			err := json.Unmarshal(value, &link)
			if err != nil {
				return err
			}

			results[string(key)] = link
			return nil
		})
		return err
	})

	return results, nil
}

// CreateLink stores a new link into database
func (db *BoltDB) CreateLink(link models.Link) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(linksBucket))

		encodedLink, err := json.Marshal(link)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(link.Name), encodedLink)
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
