package config

import (
	"encoding/json"
	"os"
)

func Load(pth string) (*Config, err) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}

	var c Config

	dec := json.NewDecoder(r)
	err = dec.Decode(&c)

	return &c, err
}

type Config struct {
	Hosts    []string `json:"hosts"`
	Keyspace string   `json:"keyspace"`
}
