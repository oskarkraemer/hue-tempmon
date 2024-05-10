package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Running hueTempratureMonitor v0.0.1")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//Load .env
	godotenv.Load()

	//DB init
	db := InitDB("./huetemp_data/sensor_logs.db")
	defer db.Close()

	//HueAPI init
	var bridge HueBridge
	bridge.AuthToken = os.Getenv("HUE_API_KEY")
	bridge.Ip = os.Getenv("HUE_BRIDGE_IP")

	var sensorRtypesFilter = []string{"temperature", "light_level"}
	var devices []Device = getSensorsByRtypes(bridge, sensorRtypesFilter)
	fmt.Printf("Found %v devices.\n", len(devices))

	ScheduleMeasurements(bridge, devices, db, "20m")

	//REST API
	ServeAPI(db, devices)

	//sig := make(chan os.Signal)
	//signal.Notify(sig, os.Interrupt, os.Kill)
	//<-sig
}
