package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func clone(dir_path string, repo string) {
	if f, err := os.Stat(dir_path); os.IsNotExist(err) || !f.IsDir() {
		err := exec.Command(
			"git", "clone",
			"https://github.com/"+repo,
			dir_path,
		).Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func convert_env(env_arr *[]string) (new_arr []string) {
	for _, env := range *env_arr {
		new_arr = append(
			new_arr,
			strings.Replace(env, "__GPMW_HOME__", GPMW_HOME, -1),
		)
	}
	return
}

func build(dir_path string, build_arr *[]string, env_arr *[]string) {
	prev, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
	}
	defer os.Chdir(prev)

	os.Chdir(dir_path)
	for _, cmd_str := range *build_arr {
		cmds := strings.Split(cmd_str, " ")
		cmd := exec.Command(cmds[0], cmds[1:]...)
		if env_arr != nil {
			new_arr := convert_env(env_arr)
			cmd.Env = append(os.Environ(), new_arr...)
		}
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func install_bin(dir_path string, cmd *string) {
	bin_path := filepath.Join(dir_path, *cmd)
	if err := os.Chmod(bin_path, 0755); err != nil {
		fmt.Println(err)
	}
	if err := os.Rename(
		bin_path,
		filepath.Join(GPMW_HOME, "bin", *cmd),
	); err != nil {
		fmt.Println(err)
	}
}

func install_src(
	src SrcPkg,
	isFin chan bool,
	wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(src.Repo)
	if src.OnOs != nil {
		if *(src.OnOs) != runtime.GOOS {
			return
		}
	}
	dir_path := filepath.Join(GPMW_HOME, "src", "github.com", src.Repo)

	if src.DoClone == nil || *(src.DoClone) {
		clone(dir_path, src.Repo)
	} else {
		InitDir(dir_path)
	}

	if src.Build != nil {
		build(dir_path, src.Build, src.BuildEnv)
	}

	if src.OnCmd != nil {
		install_bin(dir_path, src.OnCmd)
	}
	isFin <- true
}

func InstallSrc(c *SrcConfig) {
	wg := new(sync.WaitGroup)
	isFin := make(chan bool, len(c.Pkg))
	for _, src := range c.Pkg {
		wg.Add(1)
		install_src(src, isFin, wg)
	}
	wg.Wait()
	close(isFin)
	fmt.Println("Installed Packages from Source")
}
