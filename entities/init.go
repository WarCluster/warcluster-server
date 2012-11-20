package entities

type Entity interface {
    PrepareForDB() (string, []byte)
}

var Types map[string]Entity

func init() {
}

