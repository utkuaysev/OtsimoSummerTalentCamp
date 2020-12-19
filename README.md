## **OTSIMO BACKEND DEVELOPER TASK** <br />
In this task;logging will be made to a file named logfile for unexpected requests.<br />
### **CreateCandidate**<br />
`curl "http://localhost:8080/createCandidate?first_name=utku&last_name=aysev&mail=asd@gma.com&department=Development&assignee=5bb6368f55c98300013a087d&experience=true&university=TOBB"`
<br />
With this request,Firstly mail and assignee check(candidate's and assignee's department should be same) will be done.For candidate object pass this control, new candidate with __status_ **pending**_ and with _application_date_ will be inserted to database.  
### **ReadCandidate** <br />
`curl "http://localhost:8080/readCandidate?id=5b758c6151d9590001def630"`
<br />
With this request,asked candidate object will be retrieved from database and return.If id is not present in database it will be logged to file named logfile and return error message to user.
<br />
### **DeleteCandidate** <br />
`curl "http://localhost:8080/deleteCandidate?id=5fddcea1ad8dab6c97bff711""`
<br />
With this request,asked candidate object will be deleted from database.If id is not present in database it will be logged to file named logfile and return error message to user.
### **ArrangeMeeting** <br />
`curl "http://localhost:8080/arrangeMeeting?next_meeting_time=2020-12-19T09:55:33.756+00:00&id=5fddcea1ad8dab6c97bff711"`
<br />
For given _id_ and _next_meeting_time_ meeting will be arranged.There are checks for arrange meeting:<br />
_If there is a meeting not completed;new meeting can not be arranged<br />
If 4 meeting is completed;new meeting can not be arranged<br />
If last meeting is arranging; assignee will be set as Zafer_<br />
Candidate status will be set as In Progress<br />
Meeting count will increase by one.
### **CompleteMeeting** <br />
`curl "http://localhost:8080/completeMeeting?id=5fddcea1ad8dab6c97bff711""`
<br />
With this request,asked candidate's meeting will be completed and next_meeting_date field set as null date.If id is not present in database it will be logged to file named logfile and return error message to user
.**Because of requirement 2 says:If meeting count is greater than 0 and smaller than 4, the Status should be In Progress.If last meeting is completed(meeting count = 4) candidate status will be  set as Pending.**
### **DenyCandidate** <br />
`curl "http://localhost:8080/denyCandidate?id=5fddcea1ad8dab6c97bff711""`
<br />
Candidate with given id will be denied and his/her last meeting will be set null date.Status will be **Denied**
### **AcceptCandidate** <br />
`curl "http://localhost:8080/acceptCandidate?id=5fddcea1ad8dab6c97bff711""`
<br />
Candidate with given id will be checked if they completed 4 meetings.Also;*next_meeting_time* should be null date.Passed candidate will be accepted.Status will be **Accepted**
### **FindAssigneeIDByName** <br />
`curl "http://localhost:8080/findAssigneeIDByName?name=Zafer""`
<br />
If given assignee name present in database;return id.
### **FindAssigneesCandidates** <br />
`curl "http://localhost:8080/findAssigneesCandidates?id=5bb6368f55c98300013a087d""`
<br />
If given assignee id present in database;return candidates belongs to them.

