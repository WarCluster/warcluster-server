package db

import (
	"github.com/garyburd/redigo/redis"
)

// Save takes a key, name of the set and the marshaled record and writes it in.
func Save(conn redis.Conn, key, setKey string, value []byte) error {
	_, err := conn.Do("SET", key, value)
	Sadd(conn, setKey, key)
	return err
}

// Get is used to pull information from the DB in order to be used by the server.
// Get operates as read only function and does not modify the data in the DB.
func Get(conn redis.Conn, key string) ([]byte, error) {
	return redis.Bytes(conn.Do("GET", key))
}

// GetList operates as Get, but instead of an unique key it takes a patern
// in order to return a list of keys that reflect the entered patern.
func GetList(conn redis.Conn, pattern string) ([]string, error) {
	return redis.Strings(conn.Do("KEYS", pattern))
}

// Used to remove entrys from the DB.
func Delete(conn redis.Conn, key string) error {
	_, err := conn.Do("DEL", key)
	return err
}

// Saves to a redis set
func Sadd(conn redis.Conn, set, key string) error {
	_, err := conn.Do("SADD", set, key)
	return err
}

// Takes all the members in a Redis set
func Smembers(conn redis.Conn, set string) ([]interface{}, error) {
	return redis.Values(conn.Do("SMEMBERS", set))
}
