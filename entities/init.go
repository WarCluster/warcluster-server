package entities

type Entity interface {
    PrepareForDB() (string, []byte)
    GetKey() string
}

var Types map[string]Entity

