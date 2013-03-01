package entities

import (
	"../vector"
	"fmt"
)

var sunCounter = []int{0, 0}

type Sun struct {
	username string
	speed    int
	target   *vector.Vector
	Position *vector.Vector
}

func (self Sun) GetKey() string {
	return fmt.Sprintf("sun.%d_%d", self.Position.X, self.Position.Y)
}

func (self Sun) String() string {
	return fmt.Sprintf("Sun[%d, %d]", self.Position.X, self.Position.Y)
}

func (self Sun) Serialize() (string, []byte, error) {
	return self.GetKey(), []byte{1}, nil
}

func (self Sun) Update() {
	direction := self.target.Substitute(self.Position)
	if int(direction.Length()) >= self.speed {
		direction.SetLength(float64(self.speed) * ((direction.Length() / 50) + 1))
		self.Position = vector.New(float64(self.Position.X+direction.X), float64(self.Position.Y+direction.Y))
	}
}

func (self Sun) Collider(staticSun *Sun) {
	distance := self.Position.GetDistance(staticSun.Position)
	if distance < 42 { //TODO:42 da se zamesti s goleminata v pixeli na slunchevata sistema
		overlap := 42 - distance
		ndir := staticSun.Position.Substitute(self.Position)
		ndir.SetLength(overlap)
		self.Position = self.Position.Substitute(ndir)
	}
}

func (self Sun) MoveSun(position *vector.Vector) {
	self.target = position
}

func GenerateSun(friends, others []Sun) Sun {
	var newSun Sun
	var targetPosition vector.Vector

	for _, friend := range friends {
		targetPosition.X += friend.Position.X
		targetPosition.Y += friend.Position.Y
	}
	targetPosition.X /= float64(len(friends))
	targetPosition.Y /= float64(len(friends))

	noChange := false

	for noChange != true {
		var oldPos = newSun.Position
		newSun.Update()
		for _, sunEntity := range append(friends, others...) {
			newSun.Collider(&sunEntity)
		}
		if newSun.Position.IsEqual(oldPos) {
			noChange = true
		}
	}
	return newSun
	//Base player placement on worker movement from BotWars
}
