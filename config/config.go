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
	}
}

func (c *Config) Load(name string) {
	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), name)
	if err := gcfg.ReadFileInto(c, configPath); err != nil {
		log.Fatal("Error loading cfg:", err)
	}
}
