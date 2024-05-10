package main

import "time"

type RequestParmas struct {
	Table    string `uri:"table" binding:"required"`
	MinDate  string `uri:"min_date" binding:"required"`
	SensorID string `uri:"sensor_id"`
}

type SensorEntry struct {
	Id       int
	Value    float32
	SensorId string
	Created  time.Time
}
