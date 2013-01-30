package entities

import (
	"fmt"
)

var sunCounter = []int{0, 0}

type Sun []int

func (self Sun) GetKey() string {
	return fmt.Sprintf("sun.%d_%d", self[0], self[1])
}

func (self Sun) String() string {
	return fmt.Sprintf("Sun[%d, %d]", self[0], self[1])
}

func (self Sun) Serialize() (string, []byte, error) {
	return self.GetKey(), []byte{1}, nil
}

func GenerateSun() Sun {

	product := Sun{0, 0}

	product[0] = 450 + 900*sunCounter[0]
	product[1] = 450 + 900*sunCounter[1]

	if sunCounter[0] <= 1000 {
		sunCounter[0] += 1
	} else {
		sunCounter[0] = 0
		sunCounter[1] += 1
	}
	return product
	//Base player placement on worker movement from BotWars
}
