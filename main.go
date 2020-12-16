package main

import (
	"context"
	"fmt"
	"log"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
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
	// Rest of the code will go here
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//collection := client.Database("Otsimo").Collection("Candidates")
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
