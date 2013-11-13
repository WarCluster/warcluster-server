package db

import (
	"github.com/garyburd/redigo/redis"
)

// Save takes a key (struct used as template for all data containers to ease the managing of the DB)
// and generates an unique key in order to add the record to the DB.
func Save(key string, value []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

// Get is used to pull information from the DB in order to be used by the server.
// Get operates as read only function and does not modify the data in the DB.
func Get(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", key))
}

// GetList operates as Get, but instead of an unique key it takes a patern in order to return
// a list of keys that reflect the entered patern.
func GetList(pattern string) ([]interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Values(conn.Do("KEYS", pattern))
}

// I think Delete speaks for itself but still. This function is used to remove entrys from the DB.
func Delete(key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// Saves to a sorted set
func Zadd(set string, weight float64, key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZADD", set, weight, key)
	return err
}
