package random_tools

import "time"

type RandomNumberGeneration struct {
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Generated float64   `json:"generated"`
	Time      time.Time `json:"timestamp"`
}

type RandomStringGeneration struct {
	Length  int       `json:"length"`
	Charset string    `json:"charset"`
	Result  string    `json:"result"`
	Time    time.Time `json:"timestamp"`
}
