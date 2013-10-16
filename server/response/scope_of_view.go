package response

import "warcluster/entities"

type ScopeOfView struct {
	BaseResponse
	Missions map[string]*entities.Entity
	Planets  map[string]*entities.Entity
	Suns     map[string]*entities.Entity
	Entities map[string]*entities.Entity
}

func NewScopeOfView() *ScopeOfView {
	r := new(ScopeOfView)
	r.Command = "scope_of_view_result"
	return r
}
