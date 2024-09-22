package main

import (
	"123/logger"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Nums []int `json:"Nums"`
}

type Response struct {
	Res int `json:"Res"`
}

var Logger *slog.Logger

func calc(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var respBytes []byte
	req := Request{}
	resp := Response{}
	code := http.StatusInternalServerError
	var err error

	defer func() {
		if err != nil {
			Logger.Error(err.Error())
		}

		Logger.Info("calc",
			"req", req,
			"resp", resp,
			"code", code,
			"dur", time.Since(start).Milliseconds())

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBytes)
	}()

	if r.Method != http.MethodPost {
		code = http.StatusMethodNotAllowed
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		//code = http.StatusInternalServerError
		return
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		code = http.StatusBadRequest
		return
	}

	for _, n := range req.Nums {
		resp.Res += n
	}

	respBytes, err = json.Marshal(resp)
	if err != nil {
		//code = http.StatusInternalServerError
		return
	}

	code = http.StatusOK

	return
}

func main() {
	Logger = logger.New()

	adrPath := os.Getenv("ADR_PATH")
	if adrPath == "" {
		adrPath = ":8080"
	}

	http.HandleFunc("/", calc)
	http.ListenAndServe(adrPath, nil)
}
