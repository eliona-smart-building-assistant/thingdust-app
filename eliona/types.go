package eliona

import (

)



type Data struct {
	Humidity float64 `json:"humidity"`
	Occupancy bool`json:"occupied"`
	Temperature float64 `json:"temperature"`
}