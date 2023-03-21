package eliona

type Data struct {
	Humidity    float64 `json:"humidity"`
	Occupancy   int64   `json:"occupancy"`
	Temperature float64 `json:"temperature"`
}
