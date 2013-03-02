package entities

import (
	"../vector"
	"encoding/json"
	"fmt"
)

var sunCounter = []int{0, 0}

type Sun struct {
	Username string
	speed    int
	target   *vector.Vector
	position *vector.Vector
}

func (self Sun) GetKey() string {
	return fmt.Sprintf("sun.%d_%d", self.position.X, self.position.Y)
}

func (self Sun) String() string {
	return fmt.Sprintf("Sun[%d, %d]", self.position.X, self.position.Y)
}

func (self Sun) Serialize() (string, []byte, error) {
	result, err := json.Marshal(self)
	if err != nil {
		return self.GetKey(), nil, err
	}
	return self.GetKey(), result, nil
}

func (self Sun) GetPosition() *vector.Vector {
	return self.position
}

func (self Sun) Update() {
	direction := self.target.Substitute(self.position)
	if int(direction.Length()) >= self.speed {
		direction.SetLength(float64(self.speed) * ((direction.Length() / 50) + 1))
		self.position = vector.New(float64(self.position.X+direction.X), float64(self.position.Y+direction.Y))
	}
}

func (self Sun) Collider(staticSun *Sun) {
	distance := self.position.GetDistance(staticSun.position)
	if distance < 42 { //TODO:42 da se zamesti s goleminata v pixeli na slunchevata sistema
		overlap := 42 - distance
		ndir := staticSun.position.Substitute(self.position)
		ndir.SetLength(overlap)
		self.position = self.position.Substitute(ndir)
	}
}

func (self Sun) MoveSun(position *vector.Vector) {
	self.target = position
}

func GenerateSun(friends, others []Sun) Sun {
	var newSun Sun
	var targetposition vector.Vector

	for _, friend := range friends {
		targetposition.X += friend.position.X
		targetposition.Y += friend.position.Y
	}
	targetposition.X /= float64(len(friends))
	targetposition.Y /= float64(len(friends))

	noChange := false

	for noChange != true {
		var oldPos = newSun.position
		newSun.Update()
		for _, sunEntity := range append(friends, others...) {
			newSun.Collider(&sunEntity)
		}
		if newSun.position.IsEqual(oldPos) {
			noChange = true
		}
	}
	return newSun
	//Base player placement on worker movement from BotWars
}
