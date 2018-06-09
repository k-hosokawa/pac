package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func install(pkg string, isFin chan bool, wg *sync.WaitGroup) {
	// wgの数を1つ減らす（この関数が終了した時）
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

func InstallGo(gpmw_home string, c GoConfig) {
	os.Setenv("GOPATH", gpmw_home)
	os.Setenv("GOBIN", filepath.Join(gpmw_home, "bin"))

	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(c.Repos))
	for _, pkg := range c.Repos {
		wg.Add(1)
		go install(pkg, isFin, wg)
	}

	wg.Wait()
	close(isFin)
	fmt.Println("Installed Go Packages")
}
