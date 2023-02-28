package thingdust

import (

)


type Spaces map[string]Space


type Space struct {
	Humidity float64 `json:"humidity"`
	Occupancy string `json:"occupancy"`
	Temperature float64 `json:"temperature"`
}