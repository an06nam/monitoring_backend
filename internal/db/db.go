package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"server/dieselMonitoring/internal/models"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type DB struct {
	Conn *pgx.Conn
}

func ConnectDB() *DB {
	var dbConn *pgx.Conn
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	for {
		dbConn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
		if err == nil {
			log.Println("Connected to the PostgreSQL database")
			break
		}

		log.Printf("Failed to connect to PostgreSQL database: %v\n", err)
		time.Sleep(5 * time.Second)
	}

	return &DB{Conn: dbConn}
}

func (db *DB) Close() {
	db.Conn.Close(context.Background())
}

func SaveDataToDB(db *DB, jsonPayload string) {
	var SensorData models.SensorData

	err := json.Unmarshal([]byte(jsonPayload), &SensorData)
	if err != nil {
		log.Printf("Error parsing JSON data: %v\n", err)
		return
	}

	_, err = db.Conn.Exec(context.Background(), "INSERT INTO sensor_data (timestamp, uid, rpm, temperature) VALUES ($1, $2, $3, $4)",
		time.Now(), SensorData.UID, SensorData.RPM, SensorData.Temperature)

	if err != nil {
		log.Printf("Error saving data: %v\n", err)
	}
}

func (db *DB) GetDataByUID(uid string) ([]models.SensorData, error) {
	var results []models.SensorData

	rows, err := db.Conn.Query(context.Background(), "SELECT uid, rpm, temperature FROM sensor_data WHERE uid = $1", uid)
	if err != nil {
		log.Printf("Error fetching data by UID: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.SensorData
		if err := rows.Scan(&data.UID, &data.RPM, &data.Temperature); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		results = append(results, data)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v\n", err)
		return nil, err
	}

	return results, nil
}
