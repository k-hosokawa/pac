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
	Pkgs []GoPkg `toml:pkgs`
}

type GoPkg struct {
	Repo        string  `toml:repo`
	OnZshSource *string `toml:onZshSource`
	OnZshComp   *string `toml:onZshComp`
}

func install_go(pkg GoPkg, isFin chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	// install
	fmt.Println(pkg.Repo)
	cmd := exec.Command("go", "get", "-u", "github.com/"+pkg.Repo)
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	if pkg.OnZshComp != nil {
		install_zsh_completion(pkg.Repo, *(pkg.OnZshComp))
	}

	if pkg.OnZshSource != nil {
		install_zsh_source(pkg.Repo, *(pkg.OnZshSource))
	}
	isFin <- true
}

func InstallGo(_ *cli.Context) {
	os.Setenv("GOPATH", PAC_HOME)
	os.Setenv("GOBIN", filepath.Join(PAC_HOME, "bin"))

	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(CONFIG.Go.Pkgs))
	for _, pkg := range CONFIG.Go.Pkgs {
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
