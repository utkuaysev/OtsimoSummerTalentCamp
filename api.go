package main

import (
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"time"
)

func ApiCreateCandidate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// Handle error
	}
	candidate := new(Candidate)
	if err := schema.NewDecoder().Decode(candidate, r.Form); err != nil {
		// Handle error
	}

	candidate.Status = "Pending"
	candidate.Application_date = time.Now()
	_, err := CreateCandidate(*candidate)
	if err != nil {
		log.Println(err)
	}
}
func handleRequests() {
	http.HandleFunc("/article", ApiCreateCandidate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
