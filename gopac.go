package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/urfave/cli"
)

type GoConfig struct {
	Repos []string `toml:repos`
}

func install_go(pkg string, isFin chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	// install
	fmt.Println(pkg)
	cmd := exec.Command("go", "get", "-u", "github.com/"+pkg)
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	isFin <- true
}

func InstallGo(_ *cli.Context) {
	os.Setenv("GOPATH", PAC_HOME)
	os.Setenv("GOBIN", filepath.Join(PAC_HOME, "bin"))

	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(CONFIG.Go.Repos))
	for _, pkg := range CONFIG.Go.Repos {
		wg.Add(1)
		go install_go(pkg, isFin, wg)
	}

	wg.Wait()
	close(isFin)
	fmt.Println("Installed Go Packages")
}

func init() {
	command := cli.Command{
		Name:  "go",
		Usage: "Install or Update Go Packages",
		Subcommands: []cli.Command{
			{
				Name:   "update",
				Usage:  "Update Installed Packages (same as install)",
				Action: InstallGo,
			},
			{
				Name:   "install",
				Usage:  "Install Packages from source",
				Action: InstallGo,
			},
		},
	}
	APP.Commands = append(APP.Commands, command)
}
