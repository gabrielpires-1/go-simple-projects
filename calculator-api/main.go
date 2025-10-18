package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

// TODO: logging

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", addHandler)
	mux.HandleFunc("POST /multiply", multiplyHandler)
	mux.HandleFunc("POST /divide", divisionHandler)
	mux.HandleFunc("POST /subtract", subtractHandler)
	mux.HandleFunc("POST /sum", sumHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://editor.swagger.io", "https://editor-next.swagger.io"},
		AllowedMethods:   []string{http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Println("Server started on port :3000")

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var p addPayload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	sum := p.N1 + p.N2

	writeJSONResponse(w, sum)
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	var p multiplyPayload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	res := p.N1 * p.N2

	writeJSONResponse(w, res)
}

func divisionHandler(w http.ResponseWriter, r *http.Request) {
	var p divisionPayload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if p.N2 == 0 {
		http.Error(w, "divisor cannot be 0", http.StatusBadRequest)
		return
	}

	res := p.N1 / p.N2

	writeJSONResponse(w, res)
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	var p addPayload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	res := p.N1 - p.N2

	writeJSONResponse(w, res)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	var p sumPayload

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

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

	writeJSONResponse(w, int(sum))
}
