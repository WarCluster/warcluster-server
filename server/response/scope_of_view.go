package response

import (
	"warcluster/entities"

	"github.com/Vladimiroff/vec2d"
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

// calculateCanvasSize is used to determine where is the viewable by client's area
func calculateCanvasSize(position *vec2d.Vector, resolution []int) (*vec2d.Vector, *vec2d.Vector) {
	topLeft := vec2d.New(
		position.X-float64(resolution[0])/2,
		position.Y+float64(resolution[1])/2,
	)

	bottomRight := vec2d.New(
		position.X+float64(resolution[0])/2,
		position.Y-float64(resolution[1])/2,
	)
	return topLeft, bottomRight
}
