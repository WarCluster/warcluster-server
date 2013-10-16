package response

import (
	"encoding/json"
	"time"
)

type Response interface {
	MakeATimestamp()
}

type baseResponse struct {
	Command   string
	Timestamp int64
}

func (r *baseResponse) MakeATimestamp() {
	r.Timestamp = time.Now().UnixNano() / 1e6
}

func Send(r Response, sender func([]byte)) error {
	r.MakeATimestamp()
	serialized, err := json.Marshal(r)
	if err == nil {
		sender(serialized)
	}
	return err
}
