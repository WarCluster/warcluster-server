package response

import (
	"warcluster/entities"
)

type BaseResponse struct {
	Command	    string
	Timestamp   int64
}

type ScopeOfView struct {
	BaseResponse
	Entities	map[string]map[string]*entities.Entity
}

type StateChange struct {
	ScopeOfView
}

type SendMission struct {
	BaseResponse
	Mission     *entities.Mission
}
