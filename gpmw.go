package main

var (
	GPMW_HOME string
)

func main() {
	c := LoadConf("./sample.toml")
	// c.Brew.MakeBrewFile("./Brewfile")
	// GetGpmwHome()
	// home := GetGpmwHome()
	// fmt.Println(home)

	GPMW_HOME = GetGpmwHome()
	// InitDir(GPMW_HOME)
	// InstallGo(c.Go)

	// fmt.Println(c.Src)
	InstallSrc(&(c.Src))
}
