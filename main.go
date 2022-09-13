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

type Student struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	// ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// To post the student details
func CreateStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	// fmt.Println("########")
	response.Header().Set("content-type", "application/json")
	var student Student
	json.NewDecoder(request.Body).Decode(&student)
	fmt.Println(student)
	collection := client.Database("student_db").Collection("student_data")
	fmt.Println("collection", collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, student)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		// response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

// To fetch the student data
func GetStudentEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
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

// func DeleteStudentEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	params := mux.Vars(request)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])
// 	// var student Student
// 	collection := client.Database("student_db").Collection("student_data")
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
// 		return
// 	}
// 	// json.NewEncoder(response).Encode(student)
// 	if res.DeletedCount == 0 {
// 		fmt.Println("DeleteOne document not found", res)
// 	} else {
// 		fmt.Println("DeleteOne result:", res)
// 	}
// }

// To get the list of Students
func GetStudentsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var students []Student
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://mongo642:Altrancg123@cluster0.3ptkea0.mongodb.net/test?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println("Clinet ", client)
	router := mux.NewRouter()
	router.HandleFunc("/student", CreateStudentEndpoint).Methods("POST")
	router.HandleFunc("/students", GetStudentsEndpoint).Methods("GET")
	router.HandleFunc("/student/{id}", GetStudentEndpoint).Methods("GET")
	// To delete the student record
	// router.HandleFunc("/student/delete/{id}", DeleteStudentEndpoint).Methods("DELETE")
	http.ListenAndServe(":12345", router)
}
