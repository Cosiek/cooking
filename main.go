package main

import "cooking/m/v2/pantry"

func main() {
	db := pantry.InitDatabase()
	pantry.StartServer(db)
}
