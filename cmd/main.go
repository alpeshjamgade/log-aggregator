package main

import "log-aggregator/internal/app"

func main() {
	application := app.NewApp()
	application.Start()
}
