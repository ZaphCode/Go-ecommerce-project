package main

import (
	"fmt"

	"github.com/ZaphCode/clean-arch/config"
	//"github.com/ZaphCode/clean-arch/infrastructure/api"
)

func init() {
	config.MustLoadConfig("./config")
	config.MustLoadFirebaseConfig("./config")
}

func main() {
	// go api.InitBackgroundWorker()

	// app := api.Setup()

	// app.Listen(":" + config.Get().Api.Port)
	fmt.Println("Hello world")
}
