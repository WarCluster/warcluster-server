package entities

import (
	"errors"
)

const (
	RACE_VARIATION_CNT = 6
)

type Race struct {
	ID uint16
}

func (r *Race) Color() Color {
	red := []float32{0.89215686, 0.95490196, 0.99372549, 0.28235294, 0.827450, 0.32941176}
	green := []float32{0.031372549, 0.29411765, 0.79411765, 0.64901961, 0.000000, 0.57254902}
	blue := []float32{0.054901961, 0.058823529, 0.015686275, 0.074509804, 0.788235, 0.95882353}

	return Color{red[r.ID], green[r.ID], blue[r.ID]}
}

func (r *Race) Name() string {
	names := []string{"Hackafe", "BurgasLab", "InitLab", "VarnaLab", "Space lab", "Bio Lab"}
	return names[r.ID]
}

// Creates new player after the authentication and generates color based on the unique hash
func AssignRace(raceId uint16) (*Race, error) {
	if raceId > RACE_VARIATION_CNT {
		return nil, errors.New("Race id out of range.")
	}
	race := Race{ID: raceId}
	return &race, nil
}
