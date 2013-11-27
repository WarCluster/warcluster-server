package response

import (
	"fmt"

	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)

const (
	CANVAS_OFFSET_X = 10000
	CANVAS_OFFSET_Y = 10000
)

type ScopeOfView struct {
	baseResponse
	Missions map[string]*entities.Mission
	Planets  map[string]*entities.Planet
	Suns     map[string]*entities.Sun
}

func NewScopeOfView(position *vec2d.Vector, resolution []uint16) *ScopeOfView {
	s := new(ScopeOfView)
	s.Command = "scope_of_view_result"
	s.Missions = make(map[string]*entities.Mission)
	s.Planets = make(map[string]*entities.Planet)
	s.Suns = make(map[string]*entities.Sun)

	areas := listAreas(calculateCanvasSize(position, resolution))
	entityList := entities.GetAreasMembers(areas)
	for _, entity := range entityList {
		switch entity.(type) {
		case *entities.Mission:
			s.Missions[entity.Key()] = entity.(*entities.Mission)
		case *entities.Planet:
			s.Planets[entity.Key()] = entity.(*entities.Planet)
		case *entities.Sun:
			s.Suns[entity.Key()] = entity.(*entities.Sun)
		default:
		}
	}

	return s
}

// calculateCanvasSize is used to determine where is the viewable by client's area
func calculateCanvasSize(position *vec2d.Vector, resolution []uint16) (*vec2d.Vector, *vec2d.Vector) {
	topLeft := vec2d.New(
		position.X-float64(resolution[0]+CANVAS_OFFSET_X)/2,
		position.Y+float64(resolution[1]+CANVAS_OFFSET_Y)/2,
	)

	bottomRight := vec2d.New(
		position.X+float64(resolution[0]+CANVAS_OFFSET_X)/2,
		position.Y-float64(resolution[1]+CANVAS_OFFSET_Y)/2,
	)
	return topLeft, bottomRight
}

func listAreas(topLeft, bottomRight *vec2d.Vector) []string {
	topLeftX := entities.RoundCoordinateTo(topLeft.X)
	topLeftY := entities.RoundCoordinateTo(topLeft.Y)
	bottomRightX := entities.RoundCoordinateTo(bottomRight.X)
	bottomRightY := entities.RoundCoordinateTo(bottomRight.Y)

	var output []string

	for xIter := topLeftX; xIter <= bottomRightX; xIter++ {
		for yIter := topLeftY; yIter >= bottomRightY; yIter-- {
			if xIter != 0 && yIter != 0 {
				output = append(output, fmt.Sprintf("area:%v:%v", xIter, yIter))
			}
		}
	}
	return output
}
