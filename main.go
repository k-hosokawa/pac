package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	PAC_HOME = GetPacHome()
	APP      = cli.NewApp()
	CONFIG   = *LoadConf(GetPacConfig())
)

func initialize(_ *cli.Context) error {
	InitDir(PAC_HOME)
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
