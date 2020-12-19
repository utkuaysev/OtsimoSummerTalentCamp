package main

import (
	"fmt"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
	"time"
)

var loc *time.Location

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
func init() {
	fmt.Println("test")
	loc, _ = time.LoadLocation("Europe/Istanbul")
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
	candidate.Application_date = time.Now().In(loc).String()

	candidate.ID = primitive.NewObjectID().Hex()
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
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Fprintf(w, "%v", candidate)
}
func ApiDeleteCandidate(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	err := DeleteCandidate(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Fprintf(w, "Candidate with id %s is deleted", id)
}
func ApiArrangeMeeting(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	meeting_time, ok := r.URL.Query()["next_meeting_time"]

	if !ok || len(meeting_time[0]) < 1 {
		log.Println("Url Param 'next_meeting_time' is missing")
		fmt.Fprintf(w, "Url Param 'next_meeting_time' is missing")
		return
	}

	meeting_time_str := meeting_time[0]
	meeting_time_str = strings.Replace(meeting_time_str, " ", "+", 1)
	layout := "2006-01-02T15:04:05Z07:00"
	t, err := time.Parse(layout, meeting_time_str)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	err = ArrangeMeeting(id, &t)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	log.Println(err)
	fmt.Fprintf(w, "Meeting arranged for id %s", id)
}
func ApiCompleteMeeting(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	err := CompleteMeeting(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Fprintf(w, "Meeting is completed for id:%s", id)
}
func ApiDenyCandidate(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	err := DenyCandidate(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Fprintf(w, "Candidate is denied with id:%s", id)
}
func ApiAcceptCandidate(w http.ResponseWriter, r *http.Request) {
	check, id := get_id_from_query(w, r)
	if !check {
		return
	}
	err := AcceptCandidate(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	fmt.Fprintf(w, "Candidate is accepted with id:%s", id)
}
func ApiFindAssigneeIDByName(w http.ResponseWriter, r *http.Request) {
	names, ok := r.URL.Query()["name"]

	if !ok || len(names[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		fmt.Fprintf(w, "Url Param 'name' is missing")
		return
	}

	// Query()["id"] will return an array of items,
	// we only want the single item.
	name := names[0]

	id := FindAssigneeIDByName(name)
	if id == "" {
		log.Println("Assignee not found with name %s ", name)
		fmt.Fprintf(w, "Assignee not found with name %s ", name)
		return
	}
	fmt.Fprintf(w, "Candidate id is %s for name:%s", id, name)

}

func ApiFindAssigneesCandidates(writer http.ResponseWriter, request *http.Request) {
	check, id := get_id_from_query(writer, request)
	if !check {
		return
	}
	results, err := FindAssigneesCandidates(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(writer, "%s", err.Error())
		return
	}
	fmt.Fprintf(writer, "Candidates are: %v for assignee %s", results, id)

}

func handleRequests() {
	http.HandleFunc("/createCandidate", ApiCreateCandidate)
	http.HandleFunc("/readCandidate", ApiReadCandidate)
	http.HandleFunc("/deleteCandidate", ApiDeleteCandidate)
	http.HandleFunc("/arrangeMeeting", ApiArrangeMeeting)
	http.HandleFunc("/completeMeeting", ApiCompleteMeeting)
	http.HandleFunc("/denyCandidate", ApiDenyCandidate)
	http.HandleFunc("/acceptCandidate", ApiAcceptCandidate)
	http.HandleFunc("/findAssigneeIDByName", ApiFindAssigneeIDByName)
	http.HandleFunc("/findAssigneesCandidates", ApiFindAssigneesCandidates)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()
	defer end()
	defer f.Close()
}
