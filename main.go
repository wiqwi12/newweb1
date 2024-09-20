package main

import (
	"123/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type NumsRequest struct {
	Nums []int `json:"Nums"`
}

type Response struct {
	Res int `json:"Res"`
}

var Logger *slog.Logger

func calc(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "WRONG METHOD", http.StatusMethodNotAllowed)
		Logger.Log(nil, slog.LevelError, "Wrong method", time.Now(), http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "something went wrong", err)
		Logger.Log(context.Background(), slog.LevelError, "IO READALL BODY FAILURE", err)
		return
	}

	var unmarshNums NumsRequest

	err = json.Unmarshal(body, &unmarshNums)
	if err != nil {
		fmt.Fprintf(w, "empty request", http.StatusBadRequest)
		Logger.Log(context.Background(), slog.LevelError, "Unmarshall failure", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Request array: %v\n", unmarshNums.Nums)

	var sum Response
	for _, val := range unmarshNums.Nums {
		sum.Res += val
	}

	responseData, err := json.Marshal(sum)
	if err != nil {
		fmt.Fprintf(w, "something went wrong chapter 1")
		return
	}

	durr := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	logMessage := fmt.Sprintf(
		"Request Body: %v, Answer Body: %v, Request Execution Time: %v, Status: %d",
		unmarshNums.Nums, sum.Res, durr, http.StatusOK,
	)

	Logger.Log(context.Background(), slog.LevelInfo, logMessage)
	w.Write(responseData)

}

func main() {
	Logger = logger.New()

	adrPath := os.Getenv("ADR_PATH")
	if adrPath == "" {
		adrPath = ":8080"
	}

	http.HandleFunc("/", calc)
	http.ListenAndServe(":8080", nil)
}
