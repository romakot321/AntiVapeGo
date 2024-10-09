package schemas

import (
  "encoding/json"
  "log"
)

type (
  ExternalSensorDataSchema struct {
    Guid string `json:"guid" binding:"required" redis:"guid"`
    Co2 int `json:"co2" binding:"required" redis:"co2"`
    Tvoc int `json:"tvoc" binding:"required" redis:"tvoc"`
    BatteryCharge int `json:"batteryCharge" binding:"required" redis:"batteryCharge"`
  }
)

func (s ExternalSensorDataSchema) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s ExternalSensorDataSchema) UnmarshalBinary(data []byte) error {
  json.Unmarshal(data, &s)
  log.Println("data", s)
  return nil
}
