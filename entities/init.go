package entities

import "warcluster/config"

type Entity interface {
	GetKey() string
	Serialize() (string, []byte, error)
	String() string
}

var cfg config.Entities
var Types map[string]Entity
