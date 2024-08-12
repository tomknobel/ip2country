package main

import (
	"github.com/tomknobel/ip2country/cmd/api"
)

func main() {

	app := api.NewApplication()

	app.Run()
}
