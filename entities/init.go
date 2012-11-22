package entities

type Entity interface {
    Serialize() (string, []byte)
    GetKey() string
    String() string
}

var Types map[string]Entity

