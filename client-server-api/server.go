package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// CurrencyResponse represents the structure of the API response
type CurrencyResponse struct {
	USDBRL Currency `json:"USDBRL"`
}

// Currency represents the details of the currency exchange rate
type Currency struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}


const API_URL string = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type CotacaoController struct {
	ApiTimeout time.Duration
	DbTimeout time.Duration
}

func (cc *CotacaoController) GetCotacao(w http.ResponseWriter, r *http.Request) {
	dolar, err := GetDolar(cc.ApiTimeout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dolar)
}

func GetDolar(timeout time.Duration) (Currency, error) {
	client := http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", API_URL, nil)
	if err != nil {
		return Currency{}, err
	}
	resp, err := client.Do(req)	
	if err != nil {
		return Currency{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Currency{}, err
	}

	var currencyResponse CurrencyResponse
	if err := json.Unmarshal(respBytes, &currencyResponse); err != nil {
		return Currency{}, err
	}

	return currencyResponse.USDBRL, nil
}

func main() {
	cc := CotacaoController{ ApiTimeout: time.Millisecond*20, DbTimeout: time.Millisecond*10 }

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cc.GetCotacao)
	log.Fatal(http.ListenAndServe(":8080", mux))

}