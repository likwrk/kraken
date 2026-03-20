package models

import (
	"math"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/*
[
  {
    "id": "id",
    "time": 1772468514,
    "temp1": 25.69,
    "temp2": 25.56,
    "temp3": 25.88,
    "temp4": 25.88,
    "set_temp1": 7.20,
    "set_temp2": 12.00,
    "set_temp3": 13.00,
    "set_temp4": 14.00
  }
]

*/

type SensorsSetRequest struct {
	ID        string  `json:"id"`
	Timestamp int64   `json:"time"`
	Temp1     float64 `json:"temp1"`
	Temp2     float64 `json:"temp2"`
	Temp3     float64 `json:"temp3"`
	Temp4     float64 `json:"temp4"`
	SetTemp1  float64 `json:"set_temp1"`
	SetTemp2  float64 `json:"set_temp2"`
	SetTemp3  float64 `json:"set_temp3"`
	SetTemp4  float64 `json:"set_temp4"`
}

type SensorTemp struct {
	SensorID  string    `json:"sensor_id"`
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type SensorsSet struct {
	DeviceID     string       `json:"device_id"`
	Timestamp    time.Time    `json:"timestamp"`
	Values       []SensorTemp `json:"values"`
	TargetValues []SensorTemp `json:"targetValues"`
}

func (s *SensorsSetRequest) ToSensorsSet() SensorsSet {
	var timestamp time.Time
	if s.Timestamp == 0 {
		timestamp = time.Now()
	} else {
		timestamp = time.Unix(s.Timestamp, 0)
	}
	return SensorsSet{
		DeviceID: s.ID,
		Values: []SensorTemp{
			{SensorID: "temp1", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.Temp1)},
			{SensorID: "temp2", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.Temp2)},
			{SensorID: "temp3", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.Temp3)},
			{SensorID: "temp4", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.Temp4)},
		},
		TargetValues: []SensorTemp{
			{SensorID: "temp1", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.SetTemp1)},
			{SensorID: "temp2", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.SetTemp2)},
			{SensorID: "temp3", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.SetTemp3)},
			{SensorID: "temp4", DeviceID: s.ID, Timestamp: timestamp, Value: round(s.SetTemp4)},
		},
	}
}

func round(val float64) float64 {
	return math.Round(val*10) / 10
}
