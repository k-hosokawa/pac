package main

func main() {
	c := LoadConf("./sample.toml")
	c.Brew.MakeBrewFile("./Brewfile")
}
