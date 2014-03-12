package entities

import (
	"errors"
)

const (
	SUN_TEXTURE_COUNT = 1
	FRACTIONS_COUNT   = 5
)

//SetupData is a structure to hold the data from the initial character setup
//Fraction (represented by colors) is the index of the players team
//SunTextureId is the index of the home solar system sun texture
type SetupData struct {
	Fraction     uint16
	SunTextureId uint16
}

// Database key.
func (s *SetupData) Validate() error {
	if s.SunTextureId < 0 || s.SunTextureId > SUN_TEXTURE_COUNT {
		return errors.New("Sun testure index out of range.")
	}
	if s.Fraction < 0 || s.Fraction > FRACTIONS_COUNT {
		return errors.New("Fraction index out of range.")
	}
	return nil
}
