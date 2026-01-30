package main

import (
	"example.com/event-booking/config"
	"example.com/event-booking/db"
	"example.com/event-booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db.InitDB()
	server := gin.Default()

	// Route
	routes.RegisterRoutes(server)

	server.Run(":" + cfg.AppPortString())
}
