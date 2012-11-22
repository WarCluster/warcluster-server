package entities

type Entity interface {
    GetKey() string
    Serialize() (string, []byte)
    String() string
}

var Types map[string]Entity

