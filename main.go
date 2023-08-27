package main

import (
	"BasicTrade/database"
	router "BasicTrade/routers"
)

var (
	PORT = ":7070"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(PORT)
}
