package main

import (
	"log"
	"os"
)

var dbName string = "PitcherDB"

func main() {

	logger := log.New(os.Stdout, "RockPaper ", log.LstdFlags)

	app := App{
		l: logger,
	}

	app.Initialize(dbName)
	app.Run(":9090")
}
