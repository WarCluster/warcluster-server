package response

import (
	"time"
)

type BaseResponse struct {
	Command   string
	Timestamp int64
}

func (r *BaseResponse) MakeATimestamp() {
	r.Timestamp = time.Now().UnixNano() / 1e6
}
