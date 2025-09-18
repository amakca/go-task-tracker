package main

import "go-task-tracker/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
