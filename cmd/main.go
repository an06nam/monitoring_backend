package main

import (
	"log"
	"net/http"
	"os"
	"server/dieselMonitoring/internal/api"
	"server/dieselMonitoring/internal/db"
	"server/dieselMonitoring/internal/mqtt"
	"server/dieselMonitoring/internal/websocket"

	"github.com/joho/godotenv"
)

func main() {
	// Connect to the PostgreSQL database
	dbConn := db.ConnectDB()
	defer dbConn.Close()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load .env file")
	}

	// Set up MQTT client
	mqttClient := mqtt.ConnectMQTT()
	mqtt.SubscribeToMQTT(mqttClient, "my/topic", dbConn)

	// Set up WebSocket server
	ws := websocket.NewWebSocketHandler()

	// Set up HTTP server for CRUD API
	go func() {
		api.SetupCRUDRoutes(dbConn)
	}()

	// Run WebSocket server
	log.Fatal(http.ListenAndServe(os.Getenv("WEBSOCKET_PORT"), ws))
}
