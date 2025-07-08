package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Controller struct {
	UseCase *UseCase
}

func NewController(useCase *UseCase) *Controller {
	return &Controller{
		UseCase: useCase,
	}
}

func (c *Controller) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	weather, err := c.UseCase.GetWeatherFromCEP(cep)
	if err != nil {
		if errors.Is(err, ErrCepInvalid) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		} else if errors.Is(err, ErrCepNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("erro geral: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather.Current); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}