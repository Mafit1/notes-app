package main

import (
	"os"

	"github.com/Mafit1/notes-app/internal/app"
)

func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
