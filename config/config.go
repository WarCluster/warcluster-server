// Defines the structure of our configuration file
package config

import (
	"code.google.com/p/gcfg"
	"log"
	"path"
	"runtime"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Host    string
		Port    int
		Network string
	}
}

func (self *Config) Load(name string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), name)
	if err := gcfg.ReadFileInto(self, configPath); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}

type Entities struct {
	Planets struct {
		RingOffset int
		PlanetRadius int
		PlanetCount int
		PlanetHashArgs int
	}
	Suns struct {
		RandomSpawnZoneRadius int
		SolarSystemRadius float64
	}
}

func (self *Entities) Load(name string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), name)
	if err := gcfg.ReadFileInto(self, configPath); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}
