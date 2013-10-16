package response

import (
	"encoding/json"
	"time"
)

type BaseResponse struct {
	Command   string
	Timestamp int64
}

func (r *BaseResponse) MakeATimestamp() {
	r.Timestamp = time.Now().UnixNano() / 1e6
}

func (r *BaseResponse) Send(sender func([]byte) error {
	r.MakeATimestamp()
	serialized, err := json.Marshal(r)
	if err == nil {
		sender(serialized)
	}
	return err
}
