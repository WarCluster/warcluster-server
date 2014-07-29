package response

import "warcluster/entities"

type SendMissions struct {
	baseResponse
	Missions       map[string]*entities.Mission
	FailedMissions map[string]string
}

func NewSendMissions() *SendMissions {
	r := new(SendMissions)
	r.Command = "send_missions"
	r.Missions = make(map[string]*entities.Mission)
	r.FailedMissions = make(map[string]string)
	return r
}

func (m *SendMissions) Sanitize(*entities.Player) {}
