// Package config defines the structure of our configuration file
package config

import (
	"log"
	"path"
	"runtime"

	"code.google.com/p/gcfg"
)

type Config struct {
	Server struct {
		Host string
		Port uint16
	}
	Database struct {
		Host string
		Port uint16
	}
}

func (c *Config) Load(name string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), name)
	if err := gcfg.ReadFileInto(c, configPath); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}
