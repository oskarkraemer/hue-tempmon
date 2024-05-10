package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
)

var tables = []string{"temperature", "light_level"}

func InitDB(source string) *sql.DB {
	db, err := sql.Open("sqlite", source)
	if err != nil {
		log.Fatal(err)
	}

	createSQLStatement := `
		CREATE TABLE IF NOT EXISTS temperature (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			temperature FLOAT,
			sensor_id TEXT NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS light_level (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			light_level FLOAT,
			sensor_id TEXT NOT NULL,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err = db.Exec(createSQLStatement)
	if err != nil {
		log.Fatalf("Failed to created table: %v", err.Error())
	}

	return db
}

func InsertSensorEntry(db *sql.DB, table string, value float32, sensorId string) error {
	if !slices.Contains(tables, table) {
		return errors.New("nonexistent table")
	}

	_, err := db.Exec("INSERT INTO "+table+" ("+table+", sensor_id) VALUES (?, ?)", value, sensorId)
	if err != nil {
		return errors.New("Failed to insert data: " + err.Error())
	}

	return nil
}

func QuerySensorEntries(db *sql.DB, table string, min_date string, sensor_id string) ([]SensorEntry, error) {
	var rows *sql.Rows
	var err error

	if !slices.Contains(tables, table) {
		return nil, errors.New("nonexistent table")
	}

	if sensor_id == "" {
		query := fmt.Sprintf("SELECT * FROM %s WHERE date(created) >= date(?);", table)
		rows, err = db.Query(query, min_date)
	} else {
		query := fmt.Sprintf("SELECT * FROM %s WHERE date(created) >= date(?) AND sensor_id = ?;", table)
		rows, err = db.Query(query, min_date, sensor_id)
	}

	if err != nil {
		log.Fatalf("Failed to read from table: %s", err)
	}
	defer rows.Close()

	// Iterate over the rows
	var entries = []SensorEntry{}

	for rows.Next() {
		var entry SensorEntry
		err = rows.Scan(&entry.Id, &entry.Value, &entry.SensorId, &entry.Created)
		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, entry)
	}

	// Check for errors from iterating over rows
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return entries, nil
}
