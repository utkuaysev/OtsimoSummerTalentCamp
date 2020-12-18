package main

import (
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"time"
)

func get_id_from_query(w http.ResponseWriter, r *http.Request) (bool, string) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		fmt.Fprintf(w, "Url Param 'id' is missing")
		return false, ""
	}

	// Query()["id"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	return true, key
}

func ApiCreateCandidate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		err_mssg := "Error when parsing request for create candidate"
		fmt.Fprintf(w, err_mssg)
		log.Println(err_mssg)
		return
	}
	candidate := new(Candidate)
	if err := schema.NewDecoder().Decode(candidate, r.Form); err != nil {
		err_mssg := "Error when decoding request to candidate object for create candidate"
		fmt.Fprintf(w, err_mssg)
		log.Println(err_mssg)
		return
	}
	candidate.Status = "Pending"
	candidate.Application_date = time.Now()
	_, err := CreateCandidate(*candidate)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}
	fmt.Fprintf(w, "%s is successfully added ", candidate.get_name())
}
func ApiReadCandidate(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	candidate, err := ReadCandidate(id)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, "%v", candidate)
}
func handleRequests() {
	http.HandleFunc("/createCandidate", ApiCreateCandidate)
	http.HandleFunc("/readCandidate", ApiReadCandidate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
