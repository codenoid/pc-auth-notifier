package model

// AuthLog auth log model
type AuthLog struct {
	LogID     string `json:"log_id"`
	Host      string `json:"host"`
	Username  string `json:"username"`
	Status    string `json:"status"`
	Timestamp int64  `json:"ts"`
	Raw       string `json:"raw"`
	MachineID string `json:"machine_id"`
}
