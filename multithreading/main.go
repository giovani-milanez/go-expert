package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)


type CepResult string
type CepApi uint8

const (
	BRASIL_API CepApi = 0
	VIA_CEP CepApi = 1
)

func GetCEP(api CepApi, cep string) (string, error) {
	var result string

	client := http.Client{}

	url := ""
	if api == BRASIL_API {
		url = fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	} else if api == VIA_CEP {
		url = fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	} else {
		return "", fmt.Errorf("api %d not supported", api)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	resp, err := client.Do(req)	
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	result = string(body)
	
	return result, err
}


func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Argumentos invalido, uso: %s CEP", os.Args[0])
	}
	cep := os.Args[1]
	
	viaCepChan := make(chan CepResult)
	brasilApiChan := make(chan CepResult)

	// brasil api thread
	go func() {
		// time.Sleep(time.Second * 2)
		ret, err := GetCEP(BRASIL_API, cep)
		if err != nil {
			log.Fatal("falha ao buscar cep na brasil api: ", err.Error())
		}
		brasilApiChan <- CepResult(ret)
	}()

	// via cep thread
	go func() {
		// time.Sleep(time.Second * 2)
		ret, err := GetCEP(VIA_CEP, cep)
		if err != nil {
			log.Fatal("falha ao buscar cep na via cep: ", err.Error())
		}
		viaCepChan <- CepResult(ret)
	}()

	select {
		case msg := <- brasilApiChan:
			fmt.Println("API mais rapida foi BRASIL API")
			fmt.Println(msg)
		case msg := <- viaCepChan:
			fmt.Println("API mais rapida foi VIA CEP")
			fmt.Println(msg)
		case <- time.After(time.Second * 1):
			fmt.Println("Timeout!")
	}
}