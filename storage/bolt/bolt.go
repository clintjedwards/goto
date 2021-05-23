package bolt

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/goto/config"
	"github.com/clintjedwards/goto/models"
	"github.com/clintjedwards/goto/storage"
	"github.com/clintjedwards/toolkit/tkerrors"
	"github.com/rs/zerolog/log"
)

// Bolt is a representation of the bolt datastore
type Bolt struct {
	store *bolt.DB
}

// Init creates a new boltdb with given settings
func Init(config *config.BoltConfig) (Bolt, error) {
	db := Bolt{}

	store, err := bolt.Open(config.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return Bolt{}, err
	}

	// Create root bucket if not exists
	err = store.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.LinksBucket))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Bolt{}, err
	}

	db.store = store
	log.Info().Str("path", config.Path).Msg("connected to bolt db")

	return db, nil
}

// GetLink returns a link by short name
func (db *Bolt) GetLink(id string) (models.Link, error) {

	storedLink := models.Link{}

	err := db.store.View(func(tx *bolt.Tx) error {
		linksBucket := tx.Bucket([]byte(storage.LinksBucket))

		linkRaw := linksBucket.Get([]byte(id))
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
func (db *Bolt) GetAllLinks() (map[string]models.Link, error) {

	results := map[string]models.Link{}

	db.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(storage.LinksBucket))

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
func (db *Bolt) CreateLink(link *models.Link) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(storage.LinksBucket))

		exists := bucket.Get([]byte(link.ID))
		if exists != nil {
			return tkerrors.ErrEntityExists
		}

		encodedLink, err := json.Marshal(link)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(link.ID), encodedLink)
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
func (db *Bolt) BumpHitCount(id string) error {
	storedLink := models.Link{}

	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(storage.LinksBucket))

		linkRaw := bucket.Get([]byte(id))
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

		err = bucket.Put([]byte(id), encodedLink)
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
func (db *Bolt) DeleteLink(id string) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(storage.LinksBucket))

		err := bucket.Delete([]byte(id))
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
