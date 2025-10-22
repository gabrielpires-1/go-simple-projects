package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", addHandler)
	mux.HandleFunc("POST /multiply", multiplyHandler)
	mux.HandleFunc("POST /divide", divisionHandler)
	mux.HandleFunc("POST /subtract", subtractHandler)
	mux.HandleFunc("POST /sum", sumHandler)
	mux.HandleFunc("POST /sum_matrix", sumMatrixHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://editor.swagger.io", "https://editor-next.swagger.io"},
		AllowedMethods:   []string{http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(mux)

	finalHandler := loggingMiddleware(corsHandler)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	logger.Info("Server started on port :3000")

	if err := http.ListenAndServe(":3000", finalHandler); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var p addPayload

	p, ok := decodeJSON[addPayload](w, r)
	if !ok {
		return
	}

	sum := p.N1 + p.N2

	writeJSONResponseSuccess(w, sum)
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	var p multiplyPayload

	p, ok := decodeJSON[multiplyPayload](w, r)
	if !ok {
		return
	}

	res := p.N1 * p.N2

	writeJSONResponseSuccess(w, res)
}

func divisionHandler(w http.ResponseWriter, r *http.Request) {
	var p divisionPayload

	p, ok := decodeJSON[divisionPayload](w, r)
	if !ok {
		return
	}

	if p.N2 == 0 {
		// TODO: implement request id
		err := ErrorResponse{RequestID: "123", Message: "Divisor cannot be 0.", Code: "DIVISION_BY_ZERO"}
		writeJSONResquestFailed(w, err, 400)
		return
	}

	res := p.N1 / p.N2

	writeJSONResponseSuccess(w, res)
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	var p addPayload

	p, ok := decodeJSON[addPayload](w, r)
	if !ok {
		return
	}

	res := p.N1 - p.N2

	writeJSONResponseSuccess(w, res)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	var p sumPayload

	p, ok := decodeJSON[sumPayload](w, r)
	if !ok {
		return
	}

	var sum int64
	sum = 0
	sumStr := ""
	for i := 0; i < len(p); i++ {
		sum = sum + p[i]
		numStr := strconv.FormatInt(p[i], 10)
		if i == len(p)-1 {
			sumStr += numStr
		} else if i == 0 {
			sumStr += numStr
		} else {
			sumStr += numStr + " + "
		}

	}

	writeJSONResponseSuccess(w, int(sum))
}

func sumMatrixHandler(w http.ResponseWriter, r *http.Request) {
	var p matrixSumPayload

	p, ok := decodeJSON[matrixSumPayload](w, r)
	if !ok {
		return
	}

	var matrixRes matrixResult

	matrixRes.Matrix = sumMatrix(p.Matrix1, p.Matrix2)

	if err := json.NewEncoder(w).Encode(matrixRes); err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func sumMatrix(matrix1 [3][3]int, matrix2 [3][3]int) [3][3]int {
	type RowResult struct {
		RowIndex int
		RowData  [3]int
	}

	var res [3][3]int
	rows := 3

	/*
		The sum of each row is executed cocurrently in a goroutine.
		The result of each sum is sent to a buffered channel.
		The main goroutine consumes (receives) the channel results to rebuild the final matrix in correct order, using RowIndex.
	*/
	resultsCh := make(chan RowResult, rows)

	for i := 0; i < rows; i++ {
		row := i

		go func(rowIdx int) {
			rowRes := RowResult{RowIndex: rowIdx}

			for j := 0; j < 3; j++ {
				rowRes.RowData[j] = matrix1[rowIdx][j] + matrix2[rowIdx][j]
			}

			resultsCh <- rowRes
		}(row)
	}

	for i := 0; i < rows; i++ {
		result := <-resultsCh

		res[result.RowIndex] = result.RowData
	}

	return res
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	defer r.Body.Close()

	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		apiErr := ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Body couldn't be decoded: " + err.Error(),
		}
		writeJSONResquestFailed(w, apiErr, http.StatusBadRequest)
		return v, false
	}

	return v, true
}
