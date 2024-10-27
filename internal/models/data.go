package models

import "time"

type MessageTime time.Time

type MqttData struct {
	ID          int     `json:"id"`
	Temperature float32 `json:"temperature"`
	RPM         uint32  `json:"rpm"`
}

type SensorData struct {
	Timestamp   MessageTime `json:"timestamp"`
	UID         string      `json:"uid"`
	RPM         int         `json:"rpm"`
	Temperature float32     `json:"temperature"`
}
