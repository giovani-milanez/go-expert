package main

import (
	"database/sql"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
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
	DB *sql.DB
}

func (cc *CotacaoController) GetCotacao(w http.ResponseWriter, r *http.Request) {
	dolar, err := GetDolar(cc.ApiTimeout)
	if err != nil {
		log.Println("Erro ao obter cotacao do dolar: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cc.SaveToDb(dolar)
	if err != nil {
		log.Println("Erro ao tentar persistir cotacao no banco de dados: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dolar)
}

func (cc *CotacaoController) SaveToDb(dolar Currency) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, cc.DbTimeout)
	defer cancel()

	_, err := cc.DB.ExecContext(ctx,
		 "INSERT INTO cotacao(high, low, bid, ask, create_date) VALUES(?, ?, ?, ?, ?)",
		dolar.High, dolar.Low, dolar.Bid, dolar.Ask, dolar.CreateDate )

	return err
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

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./app.db") // Open a connection to the SQLite database file named app.db
	if err != nil {
	 return nil, err
	}
 
	// SQL statement to create the cotacao table if it doesn't exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS cotacao (
	 id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 high TEXT,
	 low TEXT,
	 bid TEXT,
	 ask TEXT,
	 create_date TEXT
	);`
 
	_, err = db.Exec(sqlStmt)
	if err != nil {
	 return nil, err
	}
	return db, nil
 }

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("erro ao iniciar o banco de dados: ", err.Error())
	}

	cc := CotacaoController{ ApiTimeout: time.Millisecond*200, DbTimeout: time.Millisecond*10, DB: db }

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cc.GetCotacao)
	log.Fatal(http.ListenAndServe(":8080", mux))

}