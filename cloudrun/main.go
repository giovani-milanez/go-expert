package main

import (
	"net/http"
)

const apiKey = "f4cc9333ef53455f88d142242250807"

func main() {
	useCase := NewUseCase(apiKey)
	controller := NewController(useCase)

	http.HandleFunc("/clima", controller.GetWeather)
	http.ListenAndServe(":8080", nil)
}