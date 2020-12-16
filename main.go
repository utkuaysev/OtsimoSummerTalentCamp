package main

import (
	"time"
)

type Candidate struct {
	_id        string //Unique hash that identifies candidate.
	first_name string //First name of the candidate.
	last_name  string //Last name of the candidate.
	email      string //Contact mail of candidate.
	department string //Department that candidate applied.
	/*
		 Values are
		-Marketing
		-Design
		-Development
	*/
	university string //University of the candidate.
	experience bool   //Candidate has previous working experience or not.
	status     string //Status of the candidate.
	/*
		     Values are
			-Pending
			-In Progress
			-Denied
			-Accepted
	*/
	meeting_count int       //The order of the next meeting. The maximum meeting count is 4.
	next_meeting  time.Time //Timestamp of the next meeting between the Otsimo team and the candidate.
	assignee      string    //The id of the Otsimo team member who is responsible for this candidate.
}

func main() {

}

//func ReadCandidate (_id string) (Candidate, error){
//
//}
