package main

import (
	"errors"
	"testing"
)


func TestValidCep(t *testing.T) {
	cep := "01001000"
	result, err := GetCEP(cep)
	if err != nil {
			t.Fatalf("GetCEP(%q) returned error: %v", cep, err)
	}

	if result["localidade"] != "São Paulo" {
			t.Errorf("GetCEP(%q) = %v, want %q", cep, result["localidade"], "São Paulo")
	}
}

func TestValidCep2(t *testing.T) {
	cep := "01001-000"
	result, err := GetCEP(cep)
	if err != nil {
			t.Fatalf("GetCEP(%q) returned error: %v", cep, err)
	}

	if result["localidade"] != "São Paulo" {
			t.Errorf("GetCEP(%q) = %v, want %q", cep, result["localidade"], "São Paulo")
	}
}

func TestInvalidCep(t *testing.T) {
	cep := "adfdsfsd"
	_, err := GetCEP(cep)
	if !errors.Is(err, ErrCepInvalid) {
			t.Fatalf("GetCEP(%q) expected ErrCepInvalid, got %v", cep, err)
	}
}

func TestEmptyCep(t *testing.T) {
	cep := ""
	_, err := GetCEP(cep)
	if !errors.Is(err, ErrCepInvalid) {
			t.Fatalf("GetCEP(%q) expected ErrCepInvalid, got %v", cep, err)
	}
}

func TestNotFoundCep(t *testing.T) {
	cep := "00000000"
	_, err := GetCEP(cep)
	if !errors.Is(err, ErrCepNotFound) {
			t.Fatalf("GetCEP(%q) expected ErrCepNotFound, got %v", cep, err)
	}
}

func TestShortCep(t *testing.T) {
	cep := "0000000"
	_, err := GetCEP(cep)
	if !errors.Is(err, ErrCepInvalid) {
			t.Fatalf("GetCEP(%q) expected ErrCepInvalid, got %v", cep, err)
	}
}

func TestLongCep(t *testing.T) {
	cep := "000000000"
	_, err := GetCEP(cep)
	if !errors.Is(err, ErrCepInvalid) {
			t.Fatalf("GetCEP(%q) expected ErrCepInvalid, got %v", cep, err)
	}
}

func TestWeatherValid(t *testing.T) {
	location := "São Paulo"
	resp, err := GetTemperature(apiKey, location)
	if err != nil {
		t.Fatalf("GetTemperature(%q) returned error: %v", location, err)
	}

	if resp.Location.Name != "San Paulo" {
		t.Errorf("GetTemperature(%q) = %v, want %q", location, resp.Location.Name, "San Paulo")
	}	
}

func TestWeatherEmpty(t *testing.T) {
	location := ""
	_, err := GetTemperature(apiKey, location)
	if err == nil {
		t.Fatalf("GetTemperature(%q) expected error, got nil", location)
	}
}

func TestWeatherUnkown(t *testing.T) {
	location := "fdsfdsfsdfsd"
	_, err := GetTemperature(apiKey, location)
	if err == nil {
		t.Fatalf("GetTemperature(%q) expected error, got nil", location)
	}
}