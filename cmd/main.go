package main

import (
	"fmt"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/infrastructure/api"
	"github.com/ZaphCode/clean-arch/infrastructure/utils"
)

func init() {
	config.LoadConfig()
	config.LoadFirebaseConfig()
}

func main() {
	fmt.Println("Hello world")

	utils.PrettyPrint(config.GetConfig())

	app := api.Setup()

	app.Listen(":" + config.GetConfig().Api.Port)
}
