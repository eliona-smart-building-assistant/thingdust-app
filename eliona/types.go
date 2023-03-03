package eliona

import (

)



type Data struct {
	Humidity float64 `json:"humidity"`
	Occupancy bool `json:"occupancy"`
	Temperature float64 `json:"temperature"`
}