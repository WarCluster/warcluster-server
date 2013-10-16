package response

import (
	"warcluster/entities"
)

type ScopeOfView struct {
	baseResponse
	Missions map[string]entities.Entity
	Planets  map[string]entities.Entity
	Suns     map[string]entities.Entity
}

func NewScopeOfView() *ScopeOfView {
	s := new(ScopeOfView)
	s.Command = "scope_of_view_result"
	return s
}
