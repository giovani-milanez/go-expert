package main

import (
	"context"

	"go.opentelemetry.io/otel"
)

type UseCase struct {
	ApiKey string
}

func NewUseCase(apiKey string) *UseCase {
	return &UseCase{
		ApiKey: apiKey,
	}
}

func (uc *UseCase) GetWeatherFromCEP(cep string, ctx context.Context) (Response, error) {
	if cep == "" || len(cep) < 8 || len(cep) > 9 {
		return Response{}, ErrCepInvalid
	}
	tracer := otel.Tracer("service-b")
	ctxCep, spanCep := tracer.Start(ctx, "get-cep")
	cepDetails, err := GetCEP(cep, ctxCep)
	spanCep.End()
	if err != nil {
		return Response{}, err
	}

	ctxTemp, spanTemp := tracer.Start(ctx, "get-temperature")
	weather, err := GetTemperature(uc.ApiKey, cepDetails["localidade"], ctxTemp)
	spanTemp.End()
	if err != nil {
		return Response{}, err
	}

	return Response{
		City:  cepDetails["localidade"],
		TempC: weather.Current.TempC,
		TempF: weather.Current.TempF,
		TempK: weather.Current.TempK,
	}, nil
}