package main

type HueBridge struct {
	Ip        string
	AuthToken string
}

type HueAPIResult struct {
	Devices []Device `json:"data"`
}

type Device struct {
	Id       string          `json:"id"`
	Metadata DeviceMetadata  `json:"metadata"`
	Services []DeviceService `json:"services"`
}

type DeviceMetadata struct {
	Name string `json:"name"`
}

type DeviceService struct {
	Rid   string `json:"rid"`
	Rtype string `json:"rtype"`
}

type TemperatureServiceResult struct {
	TemperatureResults []TemperatureResult `json:"data"`
}

type TemperatureResult struct {
	ServiceID   string `json:"id"`
	Temperature struct {
		Value float32 `json:"temperature"`
	} `json:"temperature"`
}

type LightLevelServiceResult struct {
	LightLevelResults []LightLevelResult `json:"data"`
}

type LightLevelResult struct {
	ServiceID  string `json:"id"`
	LightLevel struct {
		Value float32 `json:"light_level"`
	} `json:"light"`
}
