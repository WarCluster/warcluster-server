package entities

type Entity interface {
	GetKey() string
	Serialize() (string, []byte, error)
	String() string
}

var Types map[string]Entity
