package entities

import (
	"fmt"
)

type Sun []int

func (self Sun) GetKey() string {
	return fmt.Sprintf("sun.%d_%d", self[0], self[1])
}

func (self Sun) String() string {
	return fmt.Sprintf("Sun[%d, %d]", self[0], self[1])
}

func (self Sun) Serialize() (string, []byte) {
	return self.GetKey(), []byte{1}
}

func GenerateSun() Sun {
	return Sun{500, 300}
}
