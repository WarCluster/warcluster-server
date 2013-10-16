package response

import (
	"warcluster/entities"
)

type SendMission struct {
	baseResponse
	Mission *entities.Mission
}

func NewSendMission() *SendMission {
	r := new(SendMission)
	r.Command = "send_mission"
	return r
}
