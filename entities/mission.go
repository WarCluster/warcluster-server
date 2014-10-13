package entities

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type Mission struct {
	Color      Color
	Source     embeddedPlanet
	Target     embeddedPlanet
	Type       string
	StartTime  int64
	TravelTime time.Duration // in ms.
	Player     string
	ShipCount  int32
	areaSet    string
}

// Just an internal type, used to embed source and target in Mission
type embeddedPlanet struct {
	Name     string
	Owner    string
	Position *vec2d.Vector
}

// Database key.
func (m *Mission) Key() string {
	return fmt.Sprintf("mission.%d_%s", m.StartTime, m.Source.Name)
}

// Returns the sorted set by X or Y where this entity has to be put in
func (m *Mission) AreaSet() string {
	return m.areaSet
}

// Setter of m.areaSet
func (m *Mission) SetAreaSet(value string) {
	m.areaSet = value
}

// Changes its areaset based on axis and direction and updates the db
func (m *Mission) ChangeAreaSet(axis rune, direction int8) {
	areaParts := strings.Split(m.areaSet, ":")
	x, _ := strconv.ParseInt(areaParts[1], 10, 64)
	y, _ := strconv.ParseInt(areaParts[2], 10, 64)

	if axis == 'X' {
		x += int64(direction)
		if x == 0 {
			x += int64(direction)
		}
	} else if axis == 'Y' {
		y += int64(direction)
		if y == 0 {
			y += int64(direction)
		}
	}

	oldAreaSet := m.areaSet
	m.areaSet = fmt.Sprintf("area:%d:%d", x, y)
	moveToArea(m.Key(), oldAreaSet, m.areaSet)
}

// Returns all transfer points this mission will ever cross
func (m *Mission) TransferPoints() AreaTransferPoints {
	result := make(AreaTransferPoints, 0, 10)

	fillAxises := func(startPoint, endPoint float64) (container []int64) {
		startAxis := RoundCoordinateTo(startPoint)
		endAxis := RoundCoordinateTo(endPoint)
		axises := []int64{startAxis, endAxis}
		if endAxis < startAxis {
			axises = []int64{endAxis, startAxis}
		}

		for i := axises[0] + 1; i < axises[1]; i += 1 {
			container = append(container, i*Settings.AreaSize)
		}
		return
	}

	axisDirection := func(xA, xB float64) int8 {
		if xB > xA {
			return 1
		} else if xB == xA {
			return 0
		} else {
			return -1
		}
	}

	xAxises := fillAxises(m.Source.Position.X, m.Target.Position.X)
	yAxises := fillAxises(m.Source.Position.Y, m.Target.Position.Y)

	missionVectorEquation := NewCartesianEquation(m.Source.Position, m.Target.Position)

	direction := []int8{
		axisDirection(m.Source.Position.X, m.Target.Position.X),
		axisDirection(m.Source.Position.Y, m.Target.Position.Y),
	}

	for _, axis := range xAxises {
		crossPoint := vec2d.New(float64(axis), missionVectorEquation.GetYByX(float64(axis)))
		transferPoint := &AreaTransferPoint{
			TravelTime:     calculateTravelTime(m.Source.Position, crossPoint, Settings.MissionSpeed),
			Direction:      direction[0],
			CoordinateAxis: 'X',
		}
		result = append(result, transferPoint)
	}

	for _, axis := range yAxises {
		crossPoint := vec2d.New(missionVectorEquation.GetXByY(float64(axis)), float64(axis))
		transferPoint := &AreaTransferPoint{
			TravelTime:     calculateTravelTime(m.Source.Position, crossPoint, Settings.MissionSpeed),
			Direction:      direction[1],
			CoordinateAxis: 'Y',
		}
		result = append(result, transferPoint)
	}

	sort.Sort(result)
	return result
}

// Calculates the travel time in milliseconds between two planets with given speed.
// Traveling is implemented like a simple time.Sleep from our side.
func calculateTravelTime(source, target *vec2d.Vector, speed int64) time.Duration {
	distance := vec2d.GetDistance(source, target)
	return time.Duration(distance / float64(speed) * 100)
}

// When the missionary is done traveling (a.k.a. sleeping) calls this in order
// to calculate the outcome of the battle/suppliemnt/spying on target planet.

// EndAttackMission: We have to check if the target planet is owned by the attacker.
// If that's true we simply increment the ship count on that planet. If not we do the
// math and decrease the count ship on the attacked planet. We should check if the attacker
// should own that planet, which comes with all the changing colors and owner stuff.
func (m *Mission) EndAttackMission(target *Planet) (excessShips int32, ownerHasChanged bool) {
	if target.Owner == m.Player {
		m.Target.Owner = target.Owner
		m.Type = "Supply"
		return m.EndSupplyMission(target)
	} else {
		if m.ShipCount < target.ShipCount {
			target.SetShipCount(target.ShipCount - m.ShipCount)
		} else {
			if target.IsHome {
				target.SetShipCount(0)
				excessShips = m.ShipCount - target.ShipCount
			} else {
				target.SetShipCount(m.ShipCount - target.ShipCount)
				target.Owner = m.Player
				target.Color = m.Color
				ownerHasChanged = true
			}
		}
	}
	return
}

// End Supply Mission: We simply increase the ship count and we're done :P
// If however the owner of the target planet has changed we change the mission type
// to attack.
func (m *Mission) EndSupplyMission(target *Planet) (int32, bool) {
	if target.Owner != m.Target.Owner {
		m.Type = "Attack"
		return m.EndAttackMission(target)
	}

	target.SetShipCount(target.ShipCount + m.ShipCount)
	return 0, false
}

// End Spy Mission: Create a spy report for that planet and find a way to notify the logged in
// instances of the user who sent this mission.
func (m *Mission) EndSpyMission(target *Planet) (int32, bool) {
	if target.Owner == m.Player {
		m.Target.Owner = target.Owner
		return m.EndSupplyMission(target)
	}
	CreateSpyReport(target, m)
	m.ShipCount -= 1
	return 0, false
}
