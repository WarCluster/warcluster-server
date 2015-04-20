package entities

import (
	"fmt"
	"math"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type Planet struct {
	Name                string
	Color               Color
	Position            *vec2d.Vector
	IsHome              bool
	Texture             int8
	Size                int8
	LastShipCountUpdate int64
	ShipCount           int32
	MaxShipCount        int32
	Owner               string
}

// Used only when a planet is being marshalled
type PlanetPacket struct {
	Planet
	IsSpied bool `json:",omitempty"`
}

// Database key.
func (p *Planet) Key() string {
	return fmt.Sprintf("planet.%s", p.Name)
}

// Checks if the planet has an owner or not.
func (p *Planet) HasOwner() bool {
	return len(p.Owner) > 0
}

// Returns the set by X or Y where this entity has to be put in
func (p *Planet) AreaSet() string {
	return fmt.Sprintf(
		Settings.AreaTemplate,
		RoundCoordinateTo(p.Position.X),
		RoundCoordinateTo(p.Position.Y),
	)
}

// Checks what the player could see and strips it if not
// Also updates the ship count right before marshaling
func (p *Planet) Sanitize(player *Player) *PlanetPacket {
	p.UpdateShipCount()
	packet := PlanetPacket{Planet: *p}

	if p.Owner != player.Username {
		packet.ShipCount = -1
		for _, spyReport := range player.SpyReports {
			if spyReport.Name == p.Name && spyReport.IsValid() {
				packet.ShipCount = spyReport.ShipCount
				packet.IsSpied = true
			}
		}
	}
	return &packet
}

// Returns the ship count right after the ship count update.
func (p *Planet) GetShipCount() int32 {
	p.UpdateShipCount()
	return p.ShipCount
}

// Changes the ship count right after the ship count update.
// NOTE: I'm still not sure if we need a mutex here...
func (p *Planet) SetShipCount(count int32) {
	p.UpdateShipCount()
	p.ShipCount = count
	p.LastShipCountUpdate = time.Now().Unix()
}

func ShipCountTimeMod(size int8, isHome bool) int64 {
	var timeModifier float64
	if isHome {
		timeModifier = Settings.ShipsPerMinuteHome
	} else {
		switch size {
		case 1:
			timeModifier = Settings.ShipsPerMinute1
		case 2:
			timeModifier = Settings.ShipsPerMinute2
		case 3:
			timeModifier = Settings.ShipsPerMinute3
		case 4:
			timeModifier = Settings.ShipsPerMinute4
		case 5:
			timeModifier = Settings.ShipsPerMinute5
		case 6:
			timeModifier = Settings.ShipsPerMinute6
		case 7:
			timeModifier = Settings.ShipsPerMinute7
		case 8:
			timeModifier = Settings.ShipsPerMinute8
		case 9:
			timeModifier = Settings.ShipsPerMinute9
		case 10:
			timeModifier = Settings.ShipsPerMinute10
		}
	}
	return int64(timeModifier * 10)
}

// Updates the ship count based on last time this count has
// been updated and of course the planet size.
// NOTE: If the planet is somebody's home we set a static increasion rate.
func (p *Planet) UpdateShipCount() {
	if p.HasOwner() {
		passedTime := time.Now().Unix() - p.LastShipCountUpdate
		shipDiff := int32(passedTime / ShipCountTimeMod(p.Size, p.IsHome))

		if p.ShipCount > p.MaxShipCount {
			shipDiff *= int32(Settings.ShipsDeathModifier)
			if (p.ShipCount - p.MaxShipCount) > shipDiff {
				p.ShipCount -= shipDiff
			} else {
				p.ShipCount = p.MaxShipCount
			}
		} else {
			if (p.MaxShipCount - p.ShipCount) > shipDiff {
				p.ShipCount += shipDiff
			} else {
				p.ShipCount = p.MaxShipCount
			}
		}

		p.LastShipCountUpdate = time.Now().Unix()
	}
}

// Generates all planets in a solar system, based on the user's hash.
func GeneratePlanets(nickname string, sun *Sun) ([]*Planet, *Planet) {
	hash := GenerateHash(nickname)
	hashElement := func(index int) float64 {
		return float64(hash[index]) - 48 // The offset of simbol "1" in the ascii table
	}

	result := []*Planet{}
	ringOffset := float64(Settings.PlanetsRingOffset)
	planetRadius := float64(Settings.PlanetRadius)
	homePlanetIdx := int(hashElement(Settings.PlanetCount*Settings.PlanetHashArgs + 1))

	for ix := 0; ix < Settings.PlanetCount; ix++ {
		planet := Planet{
			Color:        Color{0.78431373, 0.70588235, 0.54901961},
			Position:     new(vec2d.Vector),
			IsHome:       false,
			ShipCount:    Settings.InitialPlanetShipCount,
			MaxShipCount: 0,
			Owner:        "",
		}
		// NOTE: 5 is the distance between planets
		ringOffset += planetRadius + hashElement(4*ix)*5

		planet.Name = fmt.Sprintf("%s%v", sun.Name, ix)
		planet.Position.X = math.Floor(sun.Position.X + ringOffset*math.Cos(hashElement(4*ix+1)*40))
		planet.Position.Y = math.Floor(sun.Position.Y + ringOffset*math.Sin(hashElement(4*ix+1)*40))
		planet.Texture = int8(hashElement(4*ix + 2))
		planet.Size = 1 + int8(hashElement(4*ix+3))
		planet.LastShipCountUpdate = time.Now().Unix()
		planet.IsHome = (ix == homePlanetIdx)
		planet.MaxShipCount = int32(Settings.PlanetMaxShipsMod * (60 / ShipCountTimeMod(planet.Size, planet.IsHome))) // the 60 is added so we have the propper spm and not the mod
		// this is needed due to the last change to the ShipCountTimeMod change
		if planet.IsHome {
			planet.ShipCount = Settings.InitialHomePlanetShipCount
		}
		result = append(result, &planet)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx

	return result, result[homePlanetIdx]
}
