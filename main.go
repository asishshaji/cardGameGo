package main

var dbName string = "PitcherDB"

func main() {

	app := App{}

	app.Initialize(dbName)
	app.Run(":9090")
}
