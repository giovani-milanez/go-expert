package main

import "fmt"

var ErrCepInvalid = fmt.Errorf("cep invalido")
var ErrCepNotFound = fmt.Errorf("cep nao encontrado")

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	} `json:"current"`
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
}