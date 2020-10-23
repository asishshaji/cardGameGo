package main

import "fmt"

func main() {

	fmt.Println("Hey")
	app := App{}

	app.Initialize("test")
	app.Run(":9090")
}
