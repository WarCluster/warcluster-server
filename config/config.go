// Defines the structure of our configuration file
package config

import (
	"code.google.com/p/gcfg"
	"log"
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

func (self *Config) Load(path string) {
	if err := gcfg.ReadFileInto(self, path); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}
