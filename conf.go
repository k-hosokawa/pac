package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Brew BrewConfig `toml:brew`
	Go   GoConfig   `toml:go`
	Src  SrcConfig  `toml:src`
}

func LoadConf(path string) (config *Config) {
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		panic(err)
	}
	return
}
