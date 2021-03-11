package main

type fileIO struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`    
	Link    string `json:"link"`   
	Expiry  string `json:"expiry"` 
}