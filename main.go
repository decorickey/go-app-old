package main

import (
	"fmt"
	"go-app/app/controllers"
	"go-app/config"
	"go-app/utils"
	"log"
)

func main() {
	utils.LoggingSettings(config.Config.LogFile)
	log.Println("start app")

	var apiKey string
	fmt.Print("input:")
	fmt.Scan(&apiKey)
	if err := controllers.UpdateLatestPrograms(apiKey); err != nil {
		log.Fatalln(err)
	}

	log.Println("start web server")
	controllers.StartWebServer()
}
