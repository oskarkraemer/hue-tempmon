package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"
)

func ScheduleMeasurements(bridge HueBridge, devices []Device, db *sql.DB, interval string) {
	c := cron.New()
	c.AddFunc("0 0,20,40 * * * *", func() { DoMeasurement(bridge, devices, db) })
	c.Start()
}

func DoMeasurement(bridge HueBridge, devices []Device, db *sql.DB) {
	fmt.Printf("Doing measurements at: %v ...\n", time.Now())

	for _, device := range devices {
		fmt.Printf("DeviceID: %v | Device name: %v\n", device.Id, device.Metadata.Name)
		temperature, err := getSensorValue(bridge, device, "temperature")
		if err != nil {
			fmt.Println(err.Error())
		}

		light_level, err := getSensorValue(bridge, device, "light_level")
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("Temperature: %v | LightLevel: %v\n", temperature, light_level)
		fmt.Println("---------------------------")

		err = InsertSensorEntry(db, "temperature", temperature, device.Id)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = InsertSensorEntry(db, "light_level", light_level, device.Id)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	fmt.Println("")
}
