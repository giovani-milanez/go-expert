package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func GetCotacao(timeout time.Duration) (map[string]string, error) {
	var cotacao map[string]string

	client := http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return cotacao, err
	}
	resp, err := client.Do(req)	
	if err != nil {
		return cotacao, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return cotacao, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cotacao, err
	}

	err = json.Unmarshal(body, &cotacao)
	
	return cotacao, err
}

func SalvarCotacaoAtual(bid string) error {
	return os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", bid)), 0644)
}

func main() {

	cotacao, err := GetCotacao(time.Millisecond * 300)
	if err != nil {
		log.Fatal("Erro ao obter cotacao: ", err)
	}

	err = SalvarCotacaoAtual(cotacao["bid"])
	if err != nil {
		log.Fatal("Erro ao salvar cotacao no arquivo", err)
	}

	fmt.Println("Dólar: ", cotacao["bid"])
}