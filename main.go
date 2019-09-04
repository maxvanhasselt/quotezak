package main

import (
	"fmt"
	"log"
	"os"

	"git.code-cloppers.com/max/quotezak/app"
)

func main() {
	fmt.Println("Starting app...")
	app := app.New(os.Args)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
