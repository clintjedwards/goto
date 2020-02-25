package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/go/models"
	"github.com/clintjedwards/toolkit/tkerrors"
)

// GetLink returns a link by short name
func (db *BoltDB) GetLink(name string) (models.Link, error) {

	storedLink := models.Link{}

	err := db.store.View(func(tx *bolt.Tx) error {
		linksBucket := tx.Bucket([]byte(linksBucket))

		linkRaw := linksBucket.Get([]byte(name))
		if linkRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := json.Unmarshal(linkRaw, &storedLink)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Link{}, err
	}

	return storedLink, nil
}

// GetAllLinks returns an unpaginated list of current links
func (db *BoltDB) GetAllLinks() (map[string]models.Link, error) {

	results := map[string]models.Link{}

	db.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(linksBucket))

		err := bucket.ForEach(func(key, value []byte) error {
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

		exists := bucket.Get([]byte(link.Name))
		if exists != nil {
			return tkerrors.ErrEntityExists
		}

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

// BumpHitCount updates the hit number on a certain link
func (db *BoltDB) BumpHitCount(name string) error {
	storedLink := models.Link{}

	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(linksBucket))

		linkRaw := bucket.Get([]byte(name))
		if linkRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := json.Unmarshal(linkRaw, &storedLink)
		if err != nil {
			return err
		}

		storedLink.Hits++

		encodedLink, err := json.Marshal(storedLink)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(name), encodedLink)
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

// DeleteLink removes a link from the database
func (db *BoltDB) DeleteLink(name string) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(linksBucket))

		err := bucket.Delete([]byte(name))
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
