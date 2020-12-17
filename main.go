package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

var candidates_collection *mongo.Collection
var assignees_collection *mongo.Collection
var f *os.File
var open_file_err error

func init_collection(collection_name string) *mongo.Collection {
	return client.Database("Otsimo").Collection(collection_name)
}
func init() {
	candidates_collection = init_collection("Candidates")
	assignees_collection = init_collection("Assignees")
	f, open_file_err = os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}
func main() {
	//ReadCandidate("5b758c6151d9590001def630")
	/*	var time_check = (time.Now())
		var time_check_pointer = &time_check
		ArrangeMeeting("5b75820a51d9590001def61e", time_check_pointer)
	*/
	//	CompleteMeeting("5b75881051d9590001def62a")
	defer end()
	defer f.Close()
}

func ReadCandidate(_id string) (Candidate, error) {
	filter := bson.M{"_id": _id}
	var result Candidate
	err = candidates_collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
	}
	return result, err
}
func CreateCandidate(candidate Candidate) (Candidate, error) {
	insertResult, err := candidates_collection.InsertOne(context.TODO(), candidate)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)
	return candidate, err
}
func DeleteCandidate(_id string) error {
	_, err := candidates_collection.DeleteOne(context.TODO(), bson.M{"_id": _id})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted with id %v documents in the trainers collection\n", _id)
	return err
}
func is_next_meeting_null(_id string) bool {
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return true
	}
	return !candidate.Next_meeting.IsZero()
}
func ArrangeMeeting(_id string, nextMeetingTime *time.Time) error {
	if !is_next_meeting_null(_id) {
		err1 := fmt.Errorf("There is a meeting has not completed for id %s.You can not arrange new meeting", _id)
		log.Println(err1)
		return err1
	}
	filter := bson.D{{"_id", _id}}

	update := bson.D{
		{"$inc", bson.D{{"meeting_count", 1}}}, {"$set", bson.D{{"next_meeting", nextMeetingTime}}},
	}
	updateResult, err := candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return err
}
func CompleteMeeting(_id string) error {

	filter := bson.D{{"_id", _id}}
	update := bson.D{
		{"$set", bson.D{{"next_meeting", nil}}},
	}
	updateResult, err := candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return err
}
