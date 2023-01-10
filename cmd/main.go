package main

import (
	"fmt"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/api"
)

func init() {
	config.LoadYAMLConfig()
}

func main() {
	fmt.Println("Hello world")

	app := api.Setup()

	app.Listen(":9000")
}
