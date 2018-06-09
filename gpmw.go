package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	GPMW_HOME = GetGpmwHome()
	APP       = cli.NewApp()
	CONFIG    = *LoadConf("./sample.toml")
)

func initialize(_ *cli.Context) error {
	InitDir(GPMW_HOME)
	return nil
}

func main() {
	APP.Name = "gpmw"
	APP.Usage = "Package Manager Wrapper"
	APP.Version = "0.0.1"
	APP.Before = initialize
	APP.Setup()
	APP.Run(os.Args)
}
