package models

type AdminConfig struct {
	Config string `json:"config"`
}

type ResourceData struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
