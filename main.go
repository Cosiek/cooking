package main

import (
	"cooking/m/v2/pantry"
	"cooking/m/v2/pantry_api"
)

func main() {
	db := pantry.InitDatabase()

	go pantry.StartServer(db)
	pantry_api.StartServer()
}
