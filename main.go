package main

//In this task;logging will be made to a file named logfile for unexpected requests.

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
	f, open_file_err = os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	ceo_id = FindAssigneeIDByName(ceo_name)

}

/*
With this request,asked candidate object will be retrieved from database and return.
If id is not present in database it will be logged to file named logfile and return error message to user.
*/
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

//With this request,
//Firstly mail and assignee check(candidate's and assignee's department should be same) will be done.
//For candidate object pass this control, new candidate with __status_ **pending**_ and with _application_date_ will be inserted to database.

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

/*
With this request,asked candidate object will be deleted from database.
If id is not present in database it will be logged to file named logfile and return error message to user.
*/
func DeleteCandidate(_id string) error {
	result, err := candidates_collection.DeleteOne(context.TODO(), bson.M{"_id": _id})
	if result.DeletedCount == 0 {
		return fmt.Errorf("No record is found for id %s", _id)
	}
	fmt.Printf("Candidate with id %s is deleted\n", _id)
	return err
}

/*
For given _id_ and _next_meeting_time_ meeting will be arranged.There are checks for arrange meeting:<br />
_If there is a meeting not completed;new meeting can not be arranged<br />
If 4 meeting is completed;new meeting can not be arranged<br />
If last meeting is arranging; assignee will be set as Zafer_<br />
Candidate status will be set as In Progress<br />
Meeting count will increase by one.

*/
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
	setElements = append(setElements, bson.E{"status", "In Progress"})
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

/*
With this request,
asked candidate's meeting will be completed and next_meeting_date field set as null date.
If id is not present in database it will be logged to file named logfile and return error message to user
Because of requirement 2 says:If meeting count is greater than 0 and smaller than 4, the Status should be In Progress.
If last meeting is completed(meeting count = 4) candidate status will be set as Pending.
*/
func CompleteMeeting(_id string) error {
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return fmt.Errorf("Candidate not found for id %s", _id)
	}
	filter := bson.D{{"_id", _id}}
	var setElements bson.D
	if candidate.is_next_meeting_null() {
		return fmt.Errorf("There is no meeting to complete %s", candidate.get_name())
	}
	if candidate.is_max_number_of_meeting_reached() {
		setElements = append(setElements, bson.E{"status", "Pending"})
	}
	setElements = append(setElements, bson.E{"next_meeting", time.Time{}})

	update := bson.D{{"$set", setElements}}
	_, err = candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Meeting %d is completed for candidate: %s\n", candidate.Meeting_count, candidate.get_name())
	return err
}

/*
Candidate with given id will be denied and his/her last meeting will be set null date.
*/
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
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return err
	}
	fmt.Printf("User with name %s is denied\n", candidate.get_name())
	return err
}

/*
Candidate with given id will be checked if they completed 4 meetings.Also;*next_meeting_time* should be null date.Passed candidate will be accepted.

*/
func AcceptCandidate(_id string) error {
	candidate, err := ReadCandidate(_id)
	if err != nil {
		return err
	}
	if candidate.Meeting_count < 4 {
		return fmt.Errorf("Candidates should complete 4 meetings to be accepted.%s completed %d number of meetings", candidate.get_name(), candidate.Meeting_count)
	}
	if !candidate.Next_meeting.IsZero() {
		return fmt.Errorf("There is an meeting for %s not completed.Acceptance cancelled", candidate.get_name())
	}
	filter := bson.D{{"_id", _id}}
	update := bson.D{
		{"$set", bson.D{{"status", "Accepted"}}}}
	_, err = candidates_collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("User with name %s is accepted\n\n", candidate.get_name())
	return err
}

/*
If given assignee name present in database;return id.

*/
func FindAssigneeIDByName(name string) string {
	doc := assignees_collection.FindOne(context.TODO(), bson.M{"name": name})
	// decode user model.
	var assignee Assignee
	doc.Decode(&assignee)
	return assignee.ID
}

/*
If given assignee id present in database;return candidates belongs to them.

*/
func FindAssigneesCandidates(id string) ([]Candidate, error) {
	var results []Candidate
	cur, err := candidates_collection.Find(context.TODO(), bson.D{{"assignee", id}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {

		var elem Candidate
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
