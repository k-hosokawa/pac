package main

func main() {
	c := LoadConf("./sample.toml")
	// c.Brew.MakeBrewFile("./Brewfile")
	// GetGpmwHome()
	// home := GetGpmwHome()
	// fmt.Println(home)
	gpmw_home := GetGpmwHome()
	InitDir(gpmw_home)
	InstallGo(gpmw_home, c.Go)
}
