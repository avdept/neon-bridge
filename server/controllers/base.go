package controllers

type Alert struct {
	Message string `json:"message"`
	Level    string `json:"level"` // "warning" or "error"
}
	