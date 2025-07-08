package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestControllerGetWeather(t *testing.T) {
	useCase := NewUseCase(apiKey)
	controller := NewController(useCase)

	tests := []struct {
		name       string
		cep       string
		wantCode  int
	}{
		{"valid cep", "01001000", http.StatusOK},
		{"invalid cep", "dsff", http.StatusUnprocessableEntity},
		{"empty cep", "", http.StatusUnprocessableEntity},
		{"not found cep", "99999999", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/weather?cep="+tt.cep, nil)
			w := httptest.NewRecorder()
			controller.GetWeather(w, req)

			res := w.Result()
			if res.StatusCode != tt.wantCode {
				t.Errorf("GetWeather(%q) = %v, want %v", tt.cep, res.StatusCode, tt.wantCode)
			}
		})
	}
}