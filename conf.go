package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Brew BrewConfig `toml:brew`
	Go   GoConfig   `toml:go`
	Src  SrcConfig  `toml:src`
}

type GoConfig struct {
	Repos []string `toml:repos`
}

type SrcConfig struct {
	Pkg []SrcPkg `toml:pkg`
}

type SrcPkg struct {
	Repo     string    `toml:repo`
	DoClone  *bool     `toml:doClone`
	OnOs     *string   `toml:onOs`
	Freeze   *bool     `toml:freeze`
	Build    *[]string `toml:build`
	BuildEnv *[]string `toml:buildEnv`
	OnApp    *string   `toml:onApp`
	OnCmd    *string   `toml:onCmd`
}

func LoadConf(path string) (config *Config) {
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		panic(err)
	}
	return
}
