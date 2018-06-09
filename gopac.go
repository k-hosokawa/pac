package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

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

func InstallGo(c GoConfig) {
	os.Setenv("GOPATH", GPMW_HOME)
	os.Setenv("GOBIN", filepath.Join(GPMW_HOME, "bin"))

	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(c.Repos))
	for _, pkg := range c.Repos {
		wg.Add(1)
		go install_go(pkg, isFin, wg)
	}

	wg.Wait()
	close(isFin)
	fmt.Println("Installed Go Packages")
}
