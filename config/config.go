// Package config defines the structure of our configuration file
package config

import (
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"code.google.com/p/gcfg"
)

type Config struct {
	Server struct {
		Host    string
		Port    uint16
		Console bool
		Ticker  time.Duration
	}
	Database struct {
		Host string
		Port uint16
	}
	Twitter struct {
		ConsumerKey    string
		ConsumerSecret string
		SecureLogin    bool
	}
	Race map[string]*struct {
		Id    uint8
		Red   float32
		Green float32
		Blue  float32
	}
	Entities Entities
}

type Entities struct {
	AreaSize                   int64
	AreaTemplate               string
	InitialHomePlanetShipCount int32
	InitialPlanetShipCount     int32
	MissionSpeed               int64
	PlanetCount                int
	PlanetHashArgs             int
	PlanetRadius               uint16
	ShipsPerMinute1            float64
	ShipsPerMinute2            float64
	ShipsPerMinute3            float64
	ShipsPerMinute4            float64
	ShipsPerMinute5            float64
	ShipsPerMinute6            float64
	ShipsPerMinute7            float64
	ShipsPerMinute8            float64
	ShipsPerMinute9            float64
	ShipsPerMinute10           float64
	ShipsPerMinuteHome         float64
	ShipsDeathModifier         float64
	PlanetMaxShipsMod          int64
	PlanetsRingOffset          uint16
	SolarSystemRadius          float64
	SpyReportValidity          time.Duration
	SunCanvasOffsetX           uint64
	SunCanvasOffsetY           uint64
	SunTextures                uint16
}

var ConfigDir string

func init() {
	_, filename, _, _ := runtime.Caller(1)
	ConfigDir = path.Dir(filename)
}

func (c *Config) Load() {
	if err := gcfg.ReadFileInto(c, path.Join(ConfigDir, "config.gcfg")); err != nil {
		if os.IsNotExist(err) {
			c.LoadDefault()
		} else {
			log.Fatal("Error loading cfg:", err)
		}
	}
}

func (c *Config) LoadDefault() {
	if err := gcfg.ReadFileInto(c, path.Join(ConfigDir, "config.gcfg.default")); err != nil {
		log.Fatal("Error loading default cfg:", err)
	}
}
