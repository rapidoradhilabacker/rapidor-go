package main

// core and other packages imported here
// main is a package that contains the application's main logic
// main is not imported in other packages

import (
	"go-pulsar-websockets/src/core"
	"go-pulsar-websockets/src/pkg/ws"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // Import Swagger files
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := core.SetRapidorCustomers(); err != nil {
		log.Fatalf("Error initializing customer connections: %v", err)
	}

	env := core.LoadEnv()

	// Set up the router
	router := gin.Default()

	// Documentation endpoint : to be implemented
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	ws.SetupWebSocketRoutes(router)
	log.Println("Starting WebSocket server on port :", env.ProjectPort)
	if err := router.Run(":%s", env.ProjectPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
