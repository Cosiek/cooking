package main

import (
	"cooking/m/v2/database"
	"cooking/m/v2/pantry"
	"cooking/m/v2/pantry_api"
)

func main() {
	db := database.InitDatabase()

	go pantry.StartServer()
	pantry_api.StartServer(db)
}
