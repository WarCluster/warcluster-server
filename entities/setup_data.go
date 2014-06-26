package entities

import (
	"errors"
)

const SUN_TEXTURE_COUNT = 1

//SetupData is a structure to hold the data from the initial character setup
//Race (represented by colors) is the index of the players race
//SunTextureId is the index of the home solar system sun texture
type SetupData struct {
	Race         uint8
	SunTextureId uint16
}

// Database key.
func (s *SetupData) Validate() error {
	if s.SunTextureId < 0 || s.SunTextureId > SUN_TEXTURE_COUNT {
		return errors.New("Sun testure index out of range.")
	}
	if s.Race >= uint8(len(Races)) {
		return errors.New("Race index out of range.")
	}
	return nil
}
