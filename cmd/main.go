package main

import "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/app"

func main() {
	app := app.New()
	app.Run()
}
