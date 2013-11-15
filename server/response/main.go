// Package response defines the messages we stream to users
package response

import (
	"encoding/json"
	"time"
)

type Timestamp int64

type baseResponse struct {
	Command   string
	Timestamp Timestamp
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Now().UnixNano() / 1e6)
}

func Send(r interface{}, sender func([]byte)) error {
	serialized, err := json.Marshal(r)
	if err == nil {
		sender(serialized)
	}
	return err
}
