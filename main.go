package main

import (
	"github.com/fingo-martpedia/fingo-transaction/cmd"
	"github.com/fingo-martpedia/fingo-transaction/helpers"
)

func main() {
	helpers.SetupLogger()

	helpers.SetupDatabase()

	// go cmd.ServeGRPC()

	cmd.ServeHTTP()
}
