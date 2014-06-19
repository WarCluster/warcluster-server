package response

import (
	"fmt"
	"log"
	"github.com/Vladimiroff/vec2d"

	"warcluster/entities"
)

const (
	CANVAS_OFFSET_X = 0
	CANVAS_OFFSET_Y = 0
)

type ScopeOfView struct {
	baseResponse
	Missions     map[string]*entities.Mission
	rawPlanets   map[string]*entities.Planet
	Planets      map[string]*entities.PlanetPacket
	Suns         map[string]*entities.Sun
	CanvasPoints struct {
		TopLeft     *vec2d.Vector
		BottomRight *vec2d.Vector
	}
}

func (s *ScopeOfView) Sanitize(player *entities.Player) {
	s.Planets = SanitizePlanets(player, s.rawPlanets)
}

// TODO: Use vector for rawPlanets
func NewScopeOfView(position *vec2d.Vector, resolution []uint16) *ScopeOfView {
	topLeft, bottomRight := calculateCanvasSize(position, resolution)
	areas := listAreas(topLeft, bottomRight)
	entityList := entities.GetAreasMembers(areas)

	s := new(ScopeOfView)
	s.Command = "scope_of_view_result"
	s.Missions = make(map[string]*entities.Mission)
	s.rawPlanets = make(map[string]*entities.Planet)
	s.Planets = make(map[string]*entities.PlanetPacket)
	s.Suns = make(map[string]*entities.Sun)
	s.CanvasPoints.TopLeft = topLeft
	s.CanvasPoints.BottomRight = bottomRight

	for _, entity := range entityList {
		switch entity.(type) {
		case *entities.Mission:
			s.Missions[entity.Key()] = entity.(*entities.Mission)
		case *entities.Planet:
			s.rawPlanets[entity.Key()] = entity.(*entities.Planet)
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
	log.Println("calculateCanvasSize: ", position.X, position.Y, resolution[0], resolution[1], topLeft.X, topLeft.Y, bottomRight.X, bottomRight.Y)
	return topLeft, bottomRight
}

func listAreas(topLeft, bottomRight *vec2d.Vector) []string {
	topLeftX := entities.RoundCoordinateTo(topLeft.X)
	topLeftY := entities.RoundCoordinateTo(topLeft.Y)
	bottomRightX := entities.RoundCoordinateTo(bottomRight.X)
	bottomRightY := entities.RoundCoordinateTo(bottomRight.Y)

	log.Println("1.listAreas: ", topLeft.X, topLeft.Y, bottomRight.X, bottomRight.Y, topLeftX, topLeftY, bottomRightX, bottomRightY)

	var output []string

	for xIter := topLeftX; xIter < bottomRightX; xIter++ {
		for yIter := bottomRightY; yIter < topLeftY; yIter++ {
			if xIter != 0 && yIter != 0 {
				log.Println("### 3.listAreas: ", xIter, yIter)
				output = append(output, fmt.Sprintf("area:%v:%v", xIter, yIter))
			}
		}
	}
	return output
}
