package response

import (
	"warcluster/entities"
)

type Response struct {
	Command	    string
	Timestamp   int64
	Entities	struct{
		Missions    []*entities.Mission
		Planets     []*entities.Planet
		Suns        []*entities.Sun
	}
}

type ScopeOfView struct {
	Response
}
