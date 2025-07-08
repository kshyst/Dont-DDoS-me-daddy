package models

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Options struct {
		WindowLength        string `yaml:"windowLength"`
		AllowedRequestCount string `yaml:"allowedRequestCount"`
		RedisExpiration     string `yaml:"redisExpiration"`
		ContextTimeout      string `yaml:"contextTimeout"`
	} `yaml:"options"`
	Server struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	} `yaml:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetIntOfStrings(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Cannot convert %s to an int", s)
	}
	return i
}
