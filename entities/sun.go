package entities

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/Vladimiroff/vec2d"
	"github.com/garyburd/redigo/redis"
)

type Sun struct {
	Username string
	Name     string
	speed    int32
	target   *vec2d.Vector
	Position *vec2d.Vector
}

// Database key.
func (s *Sun) Key() string {
	return fmt.Sprintf("sun.%s", s.Name)
}

// Returns the set by X or Y where this entity has to be put in
func (s *Sun) AreaSet() string {
	return fmt.Sprintf(
		ENTITIES_AREA_TEMPLATE,
		RoundCoordinateTo(s.Position.X),
		RoundCoordinateTo(s.Position.Y),
	)
}

// Generate sun's name out of user's initials and 3-digit random number
func (s *Sun) generateName(nickname string) {
	hash, _ := strconv.ParseInt(GenerateHash(nickname)[0:18], 10, 64)
	random := rand.New(rand.NewSource(hash))
	initials := extractUsernameInitials(nickname)
	number := random.Int31n(899) + 100 // we need a 3-digit number
	s.Name = fmt.Sprintf("%s%v", initials, number)
}

func (ss *Sun) calculateAdjacentSlots() []*SolarSlot {
	verticalOffset := math.Floor((SUNS_SOLAR_SYSTEM_RADIUS * math.Sqrt(3) / 2) + 0.5)
	angeledOffsetStepX := math.Floor((SUNS_SOLAR_SYSTEM_RADIUS / 2) + 0.5)
	horizontalOffsetStepX := float64(SUNS_SOLAR_SYSTEM_RADIUS)

	slots := []*SolarSlot{
		newSolarSlot(ss.Position.X-horizontalOffsetStepX, ss.Position.Y),
		newSolarSlot(ss.Position.X+horizontalOffsetStepX, ss.Position.Y),
		newSolarSlot(ss.Position.X-angeledOffsetStepX, ss.Position.Y+verticalOffset),
		newSolarSlot(ss.Position.X-angeledOffsetStepX, ss.Position.Y-verticalOffset),
		newSolarSlot(ss.Position.X+angeledOffsetStepX, ss.Position.Y+verticalOffset),
		newSolarSlot(ss.Position.X+angeledOffsetStepX, ss.Position.Y-verticalOffset),
	}
	return slots
}

func (ss *Sun) createAdjacentSlots() {
	slots := ss.calculateAdjacentSlots()

	for _, slot := range slots {
		entity, _ := Get(slot.Key())
		if entity == nil {
			Save(slot)
		}
	}
}

// Generates the key of the start position node
func getStartSolarSlotPosition(friends []*Sun) *SolarSlot {
	targetPosition := vec2d.New(0, 0)

	verticalOffset := math.Floor(SUNS_SOLAR_SYSTEM_RADIUS * (math.Sqrt(3) / 2))

	//Find best position between all friends
	for _, friend := range friends {
		targetPosition.Collect(friend.Position)
	}
	if len(friends) > 0 {
		targetPosition.DivToFloat64(float64(len(friends)))
	}

	//math.Floor(targetPosition.Y/SUNS_SOLAR_SYSTEM_RADIUS*(math.Sqrt(3)/2) + 0.5)

	//Approximate target to nearest node
	verticalOffsetCoefficent := math.Floor((targetPosition.Y / verticalOffset) + 0.5)
	if int64(verticalOffsetCoefficent)%2 != 0 {
		targetPosition.X += SUNS_SOLAR_SYSTEM_RADIUS / 2
	}
	targetPosition.X = SUNS_SOLAR_SYSTEM_RADIUS * math.Floor((targetPosition.X/SUNS_SOLAR_SYSTEM_RADIUS)+0.5)
	targetPosition.Y = verticalOffset * verticalOffsetCoefficent
	return newSolarSlot(targetPosition.X, targetPosition.Y)

}

func findHomeSolarSlot(rootSolarSlot *SolarSlot) *SolarSlot {
	var zLevel uint32 = 1

	rootSolarSlotEntity, err := Get(rootSolarSlot.Key())
	if err != nil && err != redis.ErrNil {
		panic(err)
	}

	if rootSolarSlotEntity == nil {
		rootSolarSlot = newSolarSlot(rootSolarSlot.Position.X, rootSolarSlot.Position.Y)
	} else {
		rootSolarSlot = rootSolarSlotEntity.(*SolarSlot)
	}
	if rootSolarSlot.Data == "" {
		return rootSolarSlot
	}

	for {
		nodes := rootSolarSlot.fetchSolarSlotsLayer(zLevel)
		for _, nodeKey := range nodes {
			if node, err := Get(nodeKey); err == nil && node != nil {
				nodeEntity, _ := node.(*SolarSlot)
				if nodeEntity.Data == "" {
					return nodeEntity
				}
			}
		}
		zLevel++
	}
}

// Uses all player's twitter friends and tries to place the sun as
// close as possible to all of them. This of course could cause tons of
// overlapping. To solve this, we simply throw the sun somewhere far away
// from the desired point and start to move it to THE POINT, but carefully
// watching for collisions.
func GenerateSun(username string, friends, others []*Sun) *Sun {
	newSun := Sun{
		Username: username,
		speed:    5,
		target:   vec2d.New(0, 0),
		Position: vec2d.New(0, 0),
	}
	newSun.generateName(username)

	node := findHomeSolarSlot(getStartSolarSlotPosition(friends))
	node.Data = newSun.Key()
	Save(node)

	newSun.Position.X = node.Position.X
	newSun.Position.Y = node.Position.Y
	newSun.createAdjacentSlots()
	return &newSun
}
