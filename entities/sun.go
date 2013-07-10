package entities

import (
	"encoding/json"
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

func (self Sun) GetKey() string {
	return fmt.Sprintf("sun.%v_%v", int64(self.position.X), int64(self.position.Y))
}

func (self Sun) String() string {
	return fmt.Sprintf("Sun[%v, %v]", int64(self.position.X), int64(self.position.Y))
}

func (self Sun) Serialize() (string, []byte, error) {
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func (self Sun) GetPosition() *vec2d.Vector {
	return self.position
}

func (self *Sun) Update() {
	direction := vec2d.Sub(self.target, self.position)
	if int(direction.Length()) >= self.speed {
		direction.SetLength(float64(self.speed) * ((direction.Length() / 50) + 1))
		self.position = vec2d.New(float64(self.position.X+direction.X), float64(self.position.Y+direction.Y))
	}
}

func (self *Sun) Collider(staticSun *Sun) {
	distance := vec2d.GetDistance(self.position, staticSun.position)
	if distance < 4576 { //TODO: da se zamesti s goleminata v pixeli na slunchevata sistema
		overlap := 4576 - distance
		ndir := vec2d.Sub(staticSun.position, self.position)
		ndir.SetLength(overlap)
		self.position.Sub(ndir)
	}
}

func (self *Sun) MoveSun(position *vec2d.Vector) {
	self.target = position
}

func GenerateSun(username string, friends, others []Sun) Sun {
	newSun := Sun{username, 5, vec2d.New(0, 0), getRandomStartPosition(50000)}
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

		if int64(newSun.position.X) == int64(oldPos.X) && int64(newSun.position.Y) == int64(oldPos.Y){
			noChange = true
		}
	}
	return newSun
	//Base player placement on worker movement from BotWars
}
