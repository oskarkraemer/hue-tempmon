package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

const DISCOVER_URL = "/resource/device"
const SENSOR_URL = "/resource/sensor/"

const TEMPERATURE_URL = "/resource/temperature/"
const LIGHT_LEVEL_URL = "/resource/light_level/"

func getSensorsByRtypes(bridge HueBridge, rtypes []string) []Device {
	res := hueGet(bridge, "https://"+bridge.Ip+"/clip/v2"+DISCOVER_URL)

	var result HueAPIResult
	json.Unmarshal([]byte(res), &result)

	//Get devices whose services not matching rtype
	var matches = []Device{}

	for _, device := range result.Devices {
		for _, rtype := range rtypes {
			if isSensor(device, rtype) && !containsDevice(matches, device) {
				matches = append(matches, device)
			}
		}
	}

	return matches
}

func getSensorValue(bridge HueBridge, device Device, rtype string) (float32, error) {
	for _, service := range device.Services {
		if service.Rtype == rtype {
			switch rtype {
			case "temperature":
				res := hueGet(bridge, "https://"+bridge.Ip+"/clip/v2"+TEMPERATURE_URL+service.Rid)
				var jsonResult TemperatureServiceResult
				json.Unmarshal([]byte(res), &jsonResult)
				return jsonResult.TemperatureResults[0].Temperature.Value, nil

			case "light_level":
				res := hueGet(bridge, "https://"+bridge.Ip+"/clip/v2"+LIGHT_LEVEL_URL+service.Rid)
				var jsonResult LightLevelServiceResult
				json.Unmarshal([]byte(res), &jsonResult)
				return jsonResult.LightLevelResults[0].LightLevel.Value, nil
			}
		}
	}

	return 0.0, errors.New("Rtype: " + rtype + " not found.")
}

func isSensor(device Device, rtype string) bool {
	var res bool = false

	for _, service := range device.Services {
		res = res || service.Rtype == rtype
	}

	return res
}

func containsDevice(s []Device, d Device) bool {
	for _, device := range s {
		if device.Id == d.Id {
			return true
		}
	}

	return false
}

func hueGet(bridge HueBridge, url string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("hue-application-key", bridge.AuthToken)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalln("Failed request with status code: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
