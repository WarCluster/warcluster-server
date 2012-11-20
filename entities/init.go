package entities

type Entity interface {
    PrepareForDB() (string, []byte)
}

var Types map[string]Entity

func init() {
    Types = map[string]Entity{
        "player": Player,
        "planet": Planet,
        // "mission": new(Mission),
    }
}

