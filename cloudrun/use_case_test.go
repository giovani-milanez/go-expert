package main

import "testing"

func TestUseCaseValidCep(t *testing.T) {
	useCase := NewUseCase(apiKey)
	cep := "01001000"
	result, err := useCase.GetWeatherFromCEP(cep)
	if err != nil {
		t.Fatalf("GetCEP(%q) returned error: %v", cep, err)
	}

	if result.Location.Name != "San Paulo" {
		t.Errorf("GetCEP(%q) = %v, want %q", cep, result.Location.Name, "San Paulo")
	}
}

func TestUseCaseInvalidCep(t *testing.T) {
	useCase := NewUseCase(apiKey)
	cep := "dsff"
	_, err := useCase.GetWeatherFromCEP(cep)
	if err != ErrCepInvalid {
		t.Fatalf("GetCEP(%q) returned %v expected ErrCepInvalid", cep, err)
	}
}

func TestUseCaseEmptyCep(t *testing.T) {
	useCase := NewUseCase(apiKey)
	cep := ""
	_, err := useCase.GetWeatherFromCEP(cep)
	if err != ErrCepInvalid {
		t.Fatalf("GetCEP(%q) returned %v expected ErrCepInvalid", cep, err)
	}
}

func TestUseCaseNotFoundCep(t *testing.T) {
	useCase := NewUseCase(apiKey)
	cep := "99999999"
	_, err := useCase.GetWeatherFromCEP(cep)
	if err != ErrCepNotFound {
		t.Fatalf("GetCEP(%q) returned %v expected ErrCepNotFound", cep, err)
	}
}