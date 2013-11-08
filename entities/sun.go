package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math/rand"
	"strconv"
)

var sunCounter = []int{0, 0}

type Sun struct {
	Username string
	Name     string
	speed    int
	target   *vec2d.Vector
	Position *vec2d.Vector
}

func (s *Sun) GetKey() string {
	return fmt.Sprintf("sun.%s", s.Name)
}

func (s *Sun) String() string {
	return fmt.Sprintf("Sun %s", s.Name)
}

func (s *Sun) Update() {
	direction := vec2d.Sub(s.target, s.Position)
	if int(direction.Length()) >= s.speed {
		direction.SetLength(float64(s.speed) * ((direction.Length() / 50) + 1))
		s.Position = vec2d.New(float64(s.Position.X+direction.X), float64(s.Position.Y+direction.Y))
	}
}

func (s *Sun) Collider(staticSun *Sun) {
	distance := vec2d.GetDistance(s.Position, staticSun.Position)
	if distance < SUNS_SOLAR_SYSTEM_RADIUS {
		overlap := SUNS_SOLAR_SYSTEM_RADIUS - distance
		ndir := vec2d.Sub(staticSun.Position, s.Position)
		ndir.SetLength(overlap)
		s.Position.Sub(ndir)
	}
}

func (s *Sun) MoveSun(position *vec2d.Vector) {
	s.target = position
}

// Generate sun's name out of user's initials and 3-digit random number
func (s *Sun) generateName(nickname string) {
	hash, _ := strconv.ParseInt(generateHash(nickname), 10, 64)
	random := rand.New(rand.NewSource(hash))
	initials := extractUsernameInitials(nickname)
	number := random.Int31n(899) + 100 // we need a 3-digit number
	s.Name = fmt.Sprintf("%s%v", initials, number)
}

func GenerateSun(username string, friends, others []Sun) *Sun {
	newSun := Sun{
		Username: username,
		Name:     "",
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
		newSun.Update()
		for _, sunEntity := range append(friends, others...) {
			newSun.Collider(&sunEntity)
		}

		if int64(newSun.Position.X) == int64(oldPos.X) && int64(newSun.Position.Y) == int64(oldPos.Y) {
			noChange = true
		}
	}
	return &newSun
	//Base player placement on worker movement from BotWars
}
