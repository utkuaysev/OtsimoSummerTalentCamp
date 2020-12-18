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
var ceo_id string
var ceo_name = "Zafer"

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
	ceo_id = FindAssigneeIDByName(ceo_name)
}

//func main() {
////	ReadCandidate("5b758c6151d9590001def630")
//
//	/*		var time_check = (time.Now())
//			var time_check_pointer = &time_check
//			ArrangeMeeting("5b75820a51d9590001def61e", time_check_pointer)
//	*/
//	//err = CompleteMeeting("5b75820a51d9590001def61e")
//	//DenyCandidate("5c4ab2429b4d8d000145d833")
//	//err = AcceptCandidate("5b75820a51d9590001def61e")
//	//fmt.Println(FindAssigneeIDByName("Sercan"))
//	//	FindAssigneesCandidates("5bb6360e55c98300013a087b")
//	//log.Fatal(err)
//	defer end()
//	defer f.Close()
//}

func ReadCandidate(_id string) (Candidate, error) {
	filter := bson.M{"_id": _id}
	var result Candidate
	err = candidates_collection.FindOne(context.TODO(), filter).Decode(&result)
	if result.ID == "" {
		return result, fmt.Errorf("Not found user with id %s", _id)
	}
	return result, err
}
func ReadAssignee(_id string) (Assignee, error) {
	filter := bson.M{"_id": _id}
	var result Assignee
	err = assignees_collection.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}
func CreateCandidate(candidate Candidate) (Candidate, error) {
	if !candidate.is_true_mail_format() {
		return candidate, fmt.Errorf("Mail format is inappropriate for adding db. %s did not added", candidate.get_name())
	}
	assignee, err := ReadAssignee(candidate.Assignee)
	if err != nil {
		return candidate, err
	}
	if candidate.Department != assignee.Department {
		return candidate, fmt.Errorf("Candidate's  and his/her Assignee's department should be same.%s did not added because his department is: %s and assignee %s department is %s ", candidate.get_name(), candidate.Department, assignee.Name, assignee.Department)
	}

	insertResult, err := candidates_collection.InsertOne(context.TODO(), candidate)
	if err != nil {
		return candidate, err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return candidate, err
}
func DeleteCandidate(_id string) error {
	result, err := candidates_collection.DeleteOne(context.TODO(), bson.M{"_id": _id})
	if result.DeletedCount == 0 {
		return fmt.Errorf("No record is found for id %s", _id)
	}
	return err
}
func ArrangeMeeting(_id string, nextMeetingTime *time.Time) error {
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return err
	}
	assignee, _ := ReadAssignee(candidate.Assignee)
	assignee_name := assignee.Name
	if !candidate.is_next_meeting_null() {
		return fmt.Errorf("There is a meeting has not completed for id %s.You can not arrange new meeting", _id)
	}
	if candidate.is_max_number_of_meeting_reached() {
		return fmt.Errorf("Maximum number of meeting is reached for id %s", _id)
	}
	var setElements bson.D
	if candidate.Status == "Pending" {
		setElements = append(setElements, bson.E{"status", "In Progress"})
	}
	if candidate.is_last_meeting_arranging() {
		setElements = append(setElements, bson.E{"assignee", ceo_id})
		assignee_name = ceo_name
	}
	setElements = append(setElements, bson.E{"next_meeting", nextMeetingTime})
	update := bson.D{
		{"$inc", bson.D{{"meeting_count", 1}}}, {"$set", setElements},
	}
	filter := bson.M{"_id": _id}
	_, err = candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Meeting %d arranged for %s with %s\n", candidate.Meeting_count+1, candidate.get_name(), assignee_name)
	return err
}

func CompleteMeeting(_id string) error {

	candidate, _ := ReadCandidate(_id)
	filter := bson.D{{"_id", _id}}
	var setElements bson.D
	if candidate.is_next_meeting_null() {
		return fmt.Errorf("There is no meeting to complete for %s", candidate.get_name())
	}
	if candidate.is_max_number_of_meeting_reached() {
		setElements = append(setElements, bson.E{"status", "Pending"})
	}
	setElements = append(setElements, bson.E{"next_meeting", time.Time{}})

	update := bson.D{{"$set", setElements}}
	_, err := candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Meeting %d is completed for candidate: %s\n", candidate.Meeting_count, candidate.get_name())
	return err
}
func DenyCandidate(_id string) error {
	filter := bson.D{{"_id", _id}}
	update := bson.D{
		{"$set", bson.D{{"status", "Denied"}}},
		{"$set", bson.D{{"next_meeting", time.Time{}}}},
	}
	_, err := candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	candidate, _ := ReadCandidate(_id)
	fmt.Printf("User with name %s is denied", candidate.get_name())
	return err
}
func AcceptCandidate(_id string) error {
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return err
	}
	if candidate.Next_meeting.IsZero() && candidate.Meeting_count == 4 {
		filter := bson.D{{"_id", _id}}
		update := bson.D{
			{"$set", bson.D{{"status", "Accepted"}}},
		}
		_, err := candidates_collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		fmt.Printf("User with name %s is accepted\n", candidate.get_name())
		return err
	}
	return fmt.Errorf("There is an obstacle for accept user with id:  %s", _id)
}
func FindAssigneeIDByName(name string) string {
	doc := assignees_collection.FindOne(context.TODO(), bson.M{"name": name})
	// decode user model.
	var assignee Assignee
	doc.Decode(&assignee)
	return assignee.ID
}
func FindAssigneesCandidates(id string) ([]Candidate, error) {
	var results []Candidate

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := candidates_collection.Find(context.TODO(), bson.D{{"assignee", id}})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Candidate
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return results, nil
}
