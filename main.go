package main

import (
	"cooking/m/v2/pantry"
	"cooking/m/v2/pantry_api"
)

func main() {

	go pantry.StartServer()
	pantry_api.StartServer(db)
}
