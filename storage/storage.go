package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Bucket represents the name of a section of key/value pairs
// usually a grouping of like items
type Bucket string

const (
	// linksBucket is a container for link objects
	linksBucket Bucket = "links"
)

// BoltDB is a representation of the bolt datastore
type BoltDB struct {
	store *bolt.DB
}

// NewBoltDB creates a new boltdb with given settings
func NewBoltDB(path string) (BoltDB, error) {
	db := BoltDB{}

	store, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return BoltDB{}, err
	}

	db.store = store

	return db, nil
}

// createBuckets creates given buckets nested inside another bucket
func (db *BoltDB) createBuckets(root *bolt.Bucket, buckets ...Bucket) error {

	for _, bucket := range buckets {
		_, err := root.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("could not create bucket: %s; %v", bucket, err)
		}
	}
	return nil
}
