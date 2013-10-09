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
	Missions map[string]*entities.Entity
	Planets  map[string]*entities.Entity
	Suns     map[string]*entities.Entity
	Entities map[string]*entities.Entity
}

type StateChange struct {
	BaseResponse
	Missions map[string]*entities.Entity `json:",omitempty"`
	Planets  map[string]*entities.Entity `json:",omitempty"`
	Suns     map[string]*entities.Entity `json:",omitempty"`
	Entities map[string]*entities.Entity `json:",omitempty"`
}

type SendMission struct {
	BaseResponse
	Mission     *entities.Mission
}
