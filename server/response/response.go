package response

import (
	"warcluster/entities"
)

type Response struct {
	Command	    string
	Timestamp   int64
	Entities	map[string]map[string]*entities.Entity
}

type ScopeOfView struct {
	Response
}
