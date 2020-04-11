package redis

import (
	"encoding/json"

	"github.com/clintjedwards/goto/config"
	"github.com/clintjedwards/goto/models"
	"github.com/clintjedwards/toolkit/tkerrors"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog/log"
)

// Redis is a representation of the redis datastore
type Redis struct {
	store *redis.Client
}

// Init creates a new db connection with given settings
func Init(config *config.RedisConfig) (Redis, error) {
	db := Redis{}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.DB,
	})

	// test if connection was established
	err := client.Ping().Err()
	if err != nil {
		return Redis{}, err
	}

	db.store = client
	log.Info().Str("host", config.Host).Msg("connected toredis")

	return db, nil
}

// GetLink returns a link by short name
func (db *Redis) GetLink(id string) (models.Link, error) {

	storedLink := models.Link{}

	linkRaw, err := db.store.Get(id).Bytes()
	if err == redis.Nil {
		return models.Link{}, tkerrors.ErrEntityNotFound
	}
	if err != nil {
		return models.Link{}, err
	}

	err = json.Unmarshal(linkRaw, &storedLink)
	if err != nil {
		return models.Link{}, err
	}

	return storedLink, nil
}

// GetAllLinks returns an unpaginated list of current links
func (db *Redis) GetAllLinks() (map[string]models.Link, error) {

	results := map[string]models.Link{}

	var cursor uint64

	for {
		keys, cursor, err := db.store.Scan(cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			var storedLink models.Link

			linkRaw, err := db.store.Get(key).Bytes()
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(linkRaw, &storedLink)
			if err != nil {
				return nil, err
			}

			results[key] = storedLink
		}

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

// CreateLink stores a new link into database
func (db *Redis) CreateLink(link models.Link) error {

	encodedLink, err := json.Marshal(link)
	if err != nil {
		return err
	}

	set, err := db.store.SetNX(link.ID, encodedLink, 0).Result()
	if !set {
		return tkerrors.ErrEntityExists
	}
	if err != nil {
		return err
	}
	return nil
}

// BumpHitCount updates the hit number on a certain link
func (db *Redis) BumpHitCount(id string) error {

	err := db.store.Watch(func(tx *redis.Tx) error {

		linkRaw, err := tx.Get(id).Bytes()
		if err == redis.Nil {
			return tkerrors.ErrEntityNotFound
		}
		if err != nil {
			return err
		}

		var storedLink models.Link
		err = json.Unmarshal(linkRaw, &storedLink)
		if err != nil {
			return err
		}

		storedLink.Hits++

		encodedLink, err := json.Marshal(storedLink)
		if err != nil {
			return err
		}

		err = tx.Set(id, encodedLink, 0).Err()
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteLink removes a link from the database
func (db *Redis) DeleteLink(id string) error {

	err := db.store.Del(id).Err()
	return err
}
