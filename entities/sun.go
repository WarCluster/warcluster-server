package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math"
	"math/rand"
	"strconv"
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

// Updates the sun position while doing this nasty placing the sun
func (s *Sun) update() {
	direction := vec2d.Sub(s.target, s.Position)
	if int32(direction.Length()) >= s.speed {
		direction.SetLength(float64(s.speed) * ((direction.Length() / 50) + 1))
		s.Position = vec2d.New(s.Position.X+direction.X, s.Position.Y+direction.Y)
	}
}

// Handles the collisions while moving the sun
func (s *Sun) collider(staticSun *Sun) {
	distance := vec2d.GetDistance(s.Position, staticSun.Position)
	if distance < SUNS_SOLAR_SYSTEM_RADIUS {
		overlap := SUNS_SOLAR_SYSTEM_RADIUS - distance
		ndir := vec2d.Sub(staticSun.Position, s.Position)
		ndir.SetLength(overlap)
		s.Position.Sub(ndir)
	}
}

// Generate sun's name out of user's initials and 3-digit random number
func (s *Sun) generateName(nickname string) {
	hash, _ := strconv.ParseInt(GenerateHash(nickname)[0:18], 10, 64)
	random := rand.New(rand.NewSource(hash))
	initials := extractUsernameInitials(nickname)
	number := random.Int31n(899) + 100 // we need a 3-digit number
	s.Name = fmt.Sprintf("%s%v", initials, number)
}

// Uses all player's twitter friends and tries to place the sun as
// close as possible to all of them. This of course could cause tons of
// overlapping. To solve this, we simply throw the sun somewhere far away
// from the desired point and start to move it to THE POINT, but carefully
// watching for collisions.
func GenerateSun(username string, friends, others []Sun) *Sun {
	newSun := Sun{
		Username: username,
		speed:    5,
		target:   vec2d.New(0, 0),
		Position: getRandomStartPosition(SUNS_RANDOM_SPAWN_ZONE_RADIUS),
	}
	newSun.generateName(username)
	targetPosition := vec2d.New(0, 0)

	for _, friend := range friends {
		targetPosition.X += friend.Position.X
		targetPosition.Y += friend.Position.Y
	}
	targetPosition.X /= float64(len(friends))
	targetPosition.Y /= float64(len(friends))

	noChange := false

	var oldPos *vec2d.Vector
	for noChange != true {
		oldPos = newSun.Position
		newSun.update()
		for _, sunEntity := range append(friends, others...) {
			newSun.collider(&sunEntity)
		}

		if newSun.Position.X == oldPos.X && newSun.Position.Y == oldPos.Y {
			noChange = true
		}
	}

	newSun.Position.X = math.Floor(newSun.Position.X)
	newSun.Position.Y = math.Floor(newSun.Position.Y)
	return &newSun
}
