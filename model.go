package main

import (
	"time"
)

type Assignee struct {
	ID         string `bson:"_id"`
	Name       string `bson:"name"`
	Department string `bson:"department"`
}

type Candidate struct {
	First_name string //First name of the candidate.
	Last_name  string //Last name of the candidate.
	Email      string //Contact mail of candidate.
	Department string //Department that candidate applied.
	/*
		 Values are
		-Marketing
		-Design
		-Development
	*/
	University string //University of the candidate.
	Experience bool   //Candidate has previous working experience or not.
	Status     string //Status of the candidate.
	/*
		     Values are
			-Pending
			-In Progress
			-Denied
			-Accepted
	*/
	Meeting_count int       //The order of the next meeting. The maximum meeting count is 4.
	Next_meeting  time.Time //Timestamp of the next meeting between the Otsimo team and the candidate.
	Assignee      string    //The id of the Otsimo team member who is responsible for this candidate.
}

func (candidate Candidate) is_next_meeting_null() bool {
	return candidate.Next_meeting.IsZero()
}
func (candidate Candidate) is_max_number_of_meeting_reached() bool {
	return candidate.Meeting_count == 4
}

func (candidate Candidate) is_last_meeting_arranging() bool {
	return candidate.Meeting_count == 3
}

func (candidate Candidate) get_name() string {
	return candidate.First_name + " " + candidate.Last_name
}
