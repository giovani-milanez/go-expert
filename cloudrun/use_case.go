package main

type UseCase struct {
	ApiKey string
}

func NewUseCase(apiKey string) *UseCase {
	return &UseCase{
		ApiKey: apiKey,
	}
}

func (uc *UseCase) GetWeatherFromCEP(cep string) (WeatherResponse, error) {
	if cep == "" || len(cep) < 8 || len(cep) > 9 {
		return WeatherResponse{}, ErrCepInvalid
	}
	cepDetails, err := GetCEP(cep)
	if err != nil {
		return WeatherResponse{}, err
	}

	weather, err := GetTemperature(uc.ApiKey, cepDetails["localidade"])
	if err != nil {
		return WeatherResponse{}, err
	}

	return weather, nil
}