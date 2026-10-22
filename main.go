// Command taskly is a colorful terminal todo manager backed by a CSV file.
package main

import (
	"os"

	"todo-cli/internal/app"
	"todo-cli/internal/storage"
)

const dataFile = "tasks.csv"

func main() {
	app.New(storage.NewCSV(dataFile), os.Stdin).Run()
}
