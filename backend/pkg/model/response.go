package model

type Response struct {
	Data      any    `json:"data"`
	DataCount int    `json:"data_count"`
	Error     string `json:"error"`
}
