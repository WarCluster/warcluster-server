package db

import (
	"github.com/garyburd/redigo/redis"
)

// Save takes a key (struct used as template for all data containers to ease the managing of the DB)
// and generates an unique key in order to add the record to the DB.
func Save(key string, value []byte) error {
	defer mutex.Unlock()
	mutex.Lock()

	_, err := connection.Do("SET", key, value)
	return err
}

// Get is used to pull information from the DB in order to be used by the server.
// Get operates as read only function and does not modify the data in the DB.
func Get(key string) ([]byte, error) {
	defer mutex.Unlock()
	mutex.Lock()

	return redis.Bytes(connection.Do("GET", key))
}

// GetList operates as Get, but instead of an unique key it takes a patern in order to return
// a list of keys that reflect the entered patern.
func GetList(pattern string) ([]interface{}, error) {
	defer mutex.Unlock()
	mutex.Lock()

	return redis.Values(connection.Do("KEYS", pattern))
}

// I think Delete speaks for itself but still. This function is used to remove entrys from the DB.
func Delete(key string) error {
	defer mutex.Unlock()
	mutex.Lock()

	_, err := connection.Do("DEL", key)
	return err
}
