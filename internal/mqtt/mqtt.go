package mqtt

import (
	"encoding/json"
	"log"
	"os"
	"server/dieselMonitoring/internal/db"
	"server/dieselMonitoring/internal/models"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

var mqttClient MQTT.Client

func ConnectMQTT() MQTT.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	opts := MQTT.NewClientOptions().AddBroker(os.Getenv("MQTT_URL")).SetClientID(os.Getenv("MQTT_CLIENT_ID"))

	for {
		client := MQTT.NewClient(opts)
		token := client.Connect()
		if token.Wait() && token.Error() == nil {
			log.Println("Connected to MQTT Broker")
			mqttClient = client
			return client
		}

		log.Printf("Failed to connect to MQTT Broker: %v\n", token.Error())
		time.Sleep(5 * time.Second)
	}
}

func SubscribeToMQTT(client MQTT.Client, topic string, dbConn *db.DB) {
	client.Subscribe(topic, 1, func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("Received raw message: %s from topic: %s\n", msg.Payload(), msg.Topic())

		var SensorData models.SensorData
		if err := json.Unmarshal(msg.Payload(), &SensorData); err != nil {
			log.Printf("Error unmarshaling JSON: %v\n", err)
			return
		}

		log.Printf("Parsed message: UID=%s, RPM=%d, Temperature=%.2f\n",
			SensorData.UID, SensorData.RPM, SensorData.Temperature)

		db.SaveDataToDB(dbConn, string(msg.Payload()))
	})
}

func CheckMQTTConnection() {
	for {
		if !mqttClient.IsConnected() {
			log.Println("MQTT connection lost. Attempt to reconnect....")
			mqttClient = ConnectMQTT()
		}
		time.Sleep(10 * time.Second)
	}
}
