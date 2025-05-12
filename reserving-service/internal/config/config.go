package config

import (
	"errors"
	"os"
	"slices"
)

type Config struct {
	brokers []string
	topic   string
}

func LoadConfig() (*Config, error) {
	broker := os.Getenv("BROKERS")
	topic := os.Getenv("TOPIC")

	if anyIsEmpty(broker, topic) {
		return nil, errors.New("empty property found")
	}

	return &Config{
		brokers: []string{broker},
		topic:   topic,
	}, nil
}

func anyIsEmpty(properies ...string) bool {
	return slices.Contains(properies, "")
}

func (c *Config) Brokers() []string {
	return c.brokers
}

func (c *Config) Topic() string {
	return c.topic
}
