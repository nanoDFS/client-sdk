package config

import (
	"os"
	"sync"

	"github.com/charmbracelet/log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Chunk struct {
		Size        int64 `yaml:"size"`
		PayloadSize int64 `yaml:"payloadsize"`
	} `yaml:"Chunk"`
	Master struct {
		Addr string `yaml:"address"`
	} `yaml:"Master"`
}

var mu sync.RWMutex = sync.RWMutex{}
var config *Config

func LoadConfig() *Config {
	mu.Lock()
	defer mu.Unlock()
	if config != nil {
		return config
	}
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling yaml: %v", err)
	}
	log.Infof("Successfully loaded config: %v", config)
	return config
}
