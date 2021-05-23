package storage

import (
	"github.com/clintjedwards/goto/models"
)

// Bucket represents the name of a section of key/value pairs
// usually a grouping of some sort
// ex. A key/value pair of userid-userdata would belong in the users bucket
type Bucket string

const (
	// LinksBucket represents the container in which shortened links are managed
	LinksBucket Bucket = "links"
)

// EngineType represents the different possible storage engines available
type EngineType string

const (
	// BoltEngine represents a bolt storage engine.
	// A file based key-value store.(https://github.com/boltdb/bolt)
	BoltEngine EngineType = "bolt"
	// RedisEngine represents a redis storage engine.
	// (https://redis.io/)
	RedisEngine EngineType = "redis"
)

// Engine represents backend storage implementations where items can be persisted
type Engine interface {
	GetAllLinks() (map[string]models.Link, error)
	GetLink(id string) (models.Link, error)
	CreateLink(link *models.Link) error
	BumpHitCount(id string) error
	DeleteLink(id string) error
}
