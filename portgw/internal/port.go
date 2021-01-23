package internal

// Port is a domain type representing the a given port.
// It holds information to identify and locate the port.
type Port struct {
	Name        string        `json:"name"`
	Coordinates []float64     `json:"coordinates"`
	City        string        `json:"city"`
	Province    string        `json:"province"`
	Country     string        `json:"country"`
	Alias       []interface{} `json:"alias"`
	Regions     []interface{} `json:"regions"`
	Timezone    string        `json:"timezone"`
	Unlocs      []string      `json:"unlocs"`
	Code        string        `json:"code"`
}
