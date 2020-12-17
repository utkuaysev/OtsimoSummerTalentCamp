package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Assignee struct {
	Name       string //First name of the candidate.
	Department string //Department that candidate applied.

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

var candidates_collection *mongo.Collection
var assignees_collection *mongo.Collection

func init_collection(collection_name string) *mongo.Collection {
	return client.Database("Otsimo").Collection(collection_name)
}
func init() {
	candidates_collection = init_collection("Candidates")
	assignees_collection = init_collection("Assignees")
}
func main() {
	ReadCandidate("5b758c6151d9590001def630")
	defer end()
}

func ReadCandidate(_id string) (Candidate, error) {
	filter := bson.M{"_id": _id}
	var result Candidate
	err = candidates_collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}
