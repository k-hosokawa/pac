package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	PAC_HOME       = GetPacHome()
	APP            = cli.NewApp()
	CONFIG         = *LoadConf(GetPacConfig())
	ZSH_COMP_DIR   = filepath.Join(PAC_HOME, "etc", "zsh", "Completion")
	ZSH_SOURCE_DIR = filepath.Join(PAC_HOME, "etc", "zsh", "src")
)

func initialize(_ *cli.Context) error {
	InitDir(PAC_HOME)
	InitDir(ZSH_COMP_DIR)
	InitDir(ZSH_SOURCE_DIR)
	return nil
}

func main() {
	APP.Name = "pac"
	APP.Usage = "Package Manager Wrapper"
	APP.Version = "0.0.1"
	APP.Before = initialize
	APP.Setup()
	APP.Run(os.Args)
}
