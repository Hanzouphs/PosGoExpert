package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("Quotation.txt")
	if err != nil {
		panic(err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	bodyString := string(bodyBytes)

	_, err = file.WriteString("Dólar: " + bodyString)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File created successfully!")

}
