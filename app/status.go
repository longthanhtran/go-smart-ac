package app

type Status struct {
	AcSerial     string `json:"ac_serial"`
	Temperature  uint8  `json:"temperature"`
	Humidity     uint8  `json:"humidity"`
	CoLevel      uint16 `json:"co_level"`
	HealthStatus string `json:"health_status"`
}
