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
		Host    string
		Port    uint16
		Console bool
	}
	Database struct {
		Host string
		Port uint16
	}
	Twitter struct {
		ConsumerKey       string
		ConsumerSecret    string
		AccessToken       string
		AccessTokenSecret string
	}
	Team map[string]*struct {
		Id    uint8
		Red   float32
		Green float32
		Blue  float32
	}
}

func (c *Config) Load(name string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), name)
	if err := gcfg.ReadFileInto(c, configPath); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}
