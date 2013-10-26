package entities

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
)

var sunCounter = []int{0, 0}

type Sun struct {
	Username string
	speed    int
	target   *vec2d.Vector
	position *vec2d.Vector
}

func (s *Sun) GetKey() string {
	return fmt.Sprintf("sun.%v_%v", int64(s.position.X), int64(s.position.Y))
}

func (s *Sun) String() string {
	return fmt.Sprintf("Sun[%v, %v]", int64(s.position.X), int64(s.position.Y))
}

func (s *Sun) GetPosition() *vec2d.Vector {
	return s.position
}

func (s *Sun) Update() {
	direction := vec2d.Sub(s.target, s.position)
	if int(direction.Length()) >= s.speed {
		direction.SetLength(float64(s.speed) * ((direction.Length() / 50) + 1))
		s.position = vec2d.New(float64(s.position.X+direction.X), float64(s.position.Y+direction.Y))
	}
}

func (s *Sun) Collider(staticSun *Sun) {
	distance := vec2d.GetDistance(s.position, staticSun.position)
	if distance < SUNS_SOLAR_SYSTEM_RADIUS {
		overlap := SUNS_SOLAR_SYSTEM_RADIUS - distance
		ndir := vec2d.Sub(staticSun.position, s.position)
		ndir.SetLength(overlap)
		s.position.Sub(ndir)
	}
}

func (s *Sun) MoveSun(position *vec2d.Vector) {
	s.target = position
}

func GenerateSun(username string, friends, others []Sun) *Sun {
	newSun := Sun{username, 5, vec2d.New(0, 0), getRandomStartPosition(SUNS_RANDOM_SPAWN_ZONE_RADIUS)}
	targetposition := vec2d.New(0, 0)

	for _, friend := range friends {
		targetposition.X += friend.position.X
		targetposition.Y += friend.position.Y
	}
	targetposition.X /= float64(len(friends))
	targetposition.Y /= float64(len(friends))

	noChange := false

	var oldPos *vec2d.Vector
	for noChange != true {
		oldPos = newSun.position
		newSun.Update()
		for _, sunEntity := range append(friends, others...) {
			newSun.Collider(&sunEntity)
		}

		if int64(newSun.position.X) == int64(oldPos.X) && int64(newSun.position.Y) == int64(oldPos.Y) {
			noChange = true
		}
	}
	return &newSun
	//Base player placement on worker movement from BotWars
}
