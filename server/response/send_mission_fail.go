package response

import (
	"warcluster/entities"
)

type SendMissionFailed struct {
	baseResponse
	Error string
}

func NewSendMissionFailed() *SendMissionFailed {
	r := new(SendMissionFailed)
	r.Command = "send_mission_failed"
	return r
}

func (m *SendMissionFailed) Sanitize(*entities.Player) {}
