package response

type SendMissionFailed struct {
	baseResponse
	Error string
}

func NewSendMissionFailed() *SendMissionFailed {
	r := new(SendMissionFailed)
	r.Command = "send_mission_failed"
	return r
}
