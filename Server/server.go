package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Quotation struct {
	USDBRL struct {
		ID         int64
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
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", QuotationHandler)
	http.ListenAndServe(":8080", nil)

}

func QuotationHandler(w http.ResponseWriter, r *http.Request) {

	quotation, error := GetQuotation()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite", "file:quotation.db?cache=shared&mode=rwc")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt := `CREATE TABLE IF NOT EXISTS quotation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    codein TEXT NOT NULL,
    name TEXT NOT NULL,
    high TEXT NOT NULL,
    low TEXT NOT NULL,
    var_bid TEXT NOT NULL,
    pct_change TEXT NOT NULL,
    bid TEXT NOT NULL,
    ask TEXT NOT NULL,
    timestamp TEXT NOT NULL,
    create_date TEXT NOT NULL
	);`
	_, err = db.Exec(stmt)
	if err != nil {
		panic(err)
	}

	err = Create(db, quotation)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation.USDBRL.Bid)

}

func GetQuotation() (*Quotation, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	start := time.Now()

	req, error := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD", nil)

	if error != nil {
		panic(error)
	}

	resp, error := http.DefaultClient.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Timeout: operação levou %v e excedeu o limite de 10ms", elapsed)

	}

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		panic(error)
	}

	var q Quotation
	error = json.Unmarshal(body, &q)
	if error != nil {
		return nil, error
	}
	return &q, nil
}

func Create(db *sql.DB, quotation *Quotation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	start := time.Now()

	stmt, err := db.PrepareContext(ctx, `
		INSERT INTO quotation (
			code, codein, name, high, low, var_bid, pct_change,
			bid, ask, timestamp, create_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("Erro ao preparar statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		quotation.USDBRL.Code,
		quotation.USDBRL.Codein,
		quotation.USDBRL.Name,
		quotation.USDBRL.High,
		quotation.USDBRL.Low,
		quotation.USDBRL.VarBid,
		quotation.USDBRL.PctChange,
		quotation.USDBRL.Bid,
		quotation.USDBRL.Ask,
		quotation.USDBRL.Timestamp,
		quotation.USDBRL.CreateDate,
	)

	elapsed := time.Since(start)

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Timeout: operação levou %v e excedeu o limite de 10ms", elapsed)
		return ctx.Err()
	}

	if err != nil {
		log.Printf("Erro ao executar insert: %v", err)
		return err
	}

	log.Printf("Insert concluído em %v", elapsed)
	return nil
}
