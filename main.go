package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var clientOptions *options.ClientOptions

type Student struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	// ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// To post the student details
func CreateStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This is Insert API")
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var student Student
	client = MongoDBConnection(clientOptions)
	json.NewDecoder(request.Body).Decode(&student)
	collection := client.Database("student_db_qa").Collection("student_data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// fmt.Println("Insert data: \n", student)

	if student.Firstname == "" || student.Lastname == "" {
		// To eleminate empty record insertion
	} else {
		result, err := collection.InsertOne(ctx, student)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(result)
		return
	}
}

// To fetch the student data
func GetStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	client = MongoDBConnection(clientOptions)
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var student Student
	collection := client.Database("student_db").Collection("student_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Student{ID: id}).Decode(&student)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(student)
}

// To update the student details
func UpdateStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This is Update API")
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	var student Student
	var data = make(map[string]string)
	client = MongoDBConnection(clientOptions)
	json.NewDecoder(request.Body).Decode(&student)
	coll := client.Database("student_db").Collection("student_data")
	filter := bson.D{{"_id", student.ID}}
	var update primitive.D

	if len(student.Firstname) != 0 && len(student.Lastname) != 0 {
		update = bson.D{{"$set", bson.D{{"firstname", student.Firstname}, {"lastname", student.Lastname}}}}
	} else if len(student.Lastname) == 0 {
		update = bson.D{{"$set", bson.D{{"firstname", student.Firstname}}}}
	} else if len(student.Firstname) == 0 {
		update = bson.D{{"$set", bson.D{{"lastname", student.Lastname}}}}
	} else {
		update = nil
	}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if result.ModifiedCount == 0 {
		data["status"] = "SUCCESS"
		data["message"] = "No recrods found"
	} else {
		data["status"] = "SUCCESS"
		data["message"] = "Updated successfully"
	}
	json.NewEncoder(response).Encode(data)
}

func DeleteStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("This is Delete API")
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(request)
	var data = make(map[string]string)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	// Database connection
	client = MongoDBConnection(clientOptions)
	// var student Student
	collection := client.Database("student_db").Collection("student_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	if res.DeletedCount == 0 {
		data["status"] = "SUCCESS"
		data["message"] = "No recrods found"
	} else {
		data["status"] = "SUCCESS"
		data["message"] = "Deleted successfully"
	}
	json.NewEncoder(response).Encode(data)
}

// To get the list of Students
func GetStudentsListEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var students []Student
	// Database connection
	client = MongoDBConnection(clientOptions)
	collection := client.Database("student_db").Collection("student_data")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var student Student
		cursor.Decode(&student)
		students = append(students, student)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(students)
}

// Main function
func main() {
	fmt.Println("Starting the application...")
	clientOptions = options.Client().ApplyURI("mongodb+srv://mongo642:Altrancg123@cluster0.3ptkea0.mongodb.net/test?retryWrites=true&w=majority")
	// fmt.Println("Clinet ", client)
	router := mux.NewRouter()
	// To insert the student details
	router.HandleFunc("/student", CreateStudentEndpoint).Methods("POST", "OPTIONS")
	// To get the students list
	router.HandleFunc("/students", GetStudentsListEndpoint).Methods("GET", "OPTIONS")
	// To update the students details
	router.HandleFunc("/student/update", UpdateStudentEndpoint).Methods("PUT", "OPTIONS")
	// To fetch the student details
	router.HandleFunc("/student/{id}", GetStudentEndpoint).Methods("GET", "OPTIONS")
	// To delete the student record
	router.HandleFunc("/student/delete/{id}", DeleteStudentEndpoint).Methods("DELETE", "OPTIONS")
	http.ListenAndServe(":12345", router)

}
