package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"hospital_management_portal/database"
	"hospital_management_portal/models"
	request "hospital_management_portal/requests"
	"hospital_management_portal/response"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateDoctor(w http.ResponseWriter, r *http.Request) {

	 //let's first validate our request
	 var docCreatReq request.DoctorRequest
	 decoder := json.NewDecoder(r.Body)
	 if err := decoder.Decode(&docCreatReq); err != nil {
		 http.Error(w, "Invalid input", http.StatusBadRequest)
		 return
	 }
 
	 // Initialize the validator and validate the request data
	 validate := validator.New()
	 err := validate.Struct(docCreatReq)
	 if err != nil {
		 if errs, ok := err.(validator.ValidationErrors); ok {
			 errMessages := request.CustomeErrorMessage(errs)
			 http.Error(w, strings.Join(errMessages, ", "), http.StatusBadRequest)
			 return
		 }
		 http.Error(w, "Validation failed", http.StatusInternalServerError)
		 return
	 }
	 
	 newDoctor := models.Doctor{
		ID: uuid.NewString(), 
		FirstName: docCreatReq.FirstName,
		LastName: docCreatReq.LastName,
		Phone: docCreatReq.Phone,
		Email: docCreatReq.Email,
        CreatedAt: time.Now(), // Set created_at and updated_at
        UpdatedAt: time.Now(),
    }

	  //lets first open the collection we need to insert the user in
	  collection:= database.OpenCollection(database.Client,"doctors")
	  ctx,cancel :=context.WithTimeout(context.Background(),10*time.Second)
	  defer cancel()
	  //let's now insert the user
	  _,err=collection.InsertOne(ctx,newDoctor)
	  if err != nil {
  
		  err_response := response.ErrorResponse{
			  OK: false,
			  Message: err.Error(),
			  Status: "failed",
		  }
		  json.NewEncoder(w).Encode(err_response)
	  }

	    // Properly respond with the created user, omitting sensitive data like password
		
		w.Header().Set("Content-Type", "application/json")
	
		resResult := map[string]interface{}{
			"ok":     true,
			"status": "success",
			"message": fmt.Sprintf("Dr.%s %s added successfully", newDoctor.FirstName,newDoctor.LastName),
		}
		json.NewEncoder(w).Encode(resResult)

}

// GetUser function to handle the request
func UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Assuming `database` and `Client` are predefined in your package
    collection := database.OpenCollection(database.Client, "doctors")

	vars := mux.Vars(r)
    doctorID := vars["doctorId"] 

    // Create a filter to find the doctor by ID
    filter := bson.M{"id": doctorID}

    // Finding the user document with the given ID
    var doctor bson.M
    err := collection.FindOne(ctx, filter).Decode(&doctor)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No doctor found with given ID", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

   var res_doc [ ]bson.M 
   res_doc=append(res_doc,doctor)

    resResult := map[string]interface{}{
        "ok":     true,
        "status": "success",
        "user":   res_doc,
    }

    // Set response header
    w.Header().Set("Content-Type", "application/json")

    // Encoding the result to JSON and sending the response
    if err := json.NewEncoder(w).Encode(resResult); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

}

// GetUsers function to handle the request
func GetUsers(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel() // This ensures that context cancellation is always called

    // Assuming `database` and `Client` are predefined in your package
    collection := database.OpenCollection(database.Client, "users")

    // Finding all documents in the collection
    cur, err := collection.Find(ctx, bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cur.Close(ctx)

    // Slice to hold all decoded documents
    var users []bson.M

    // Decoding documents into users slice
    for cur.Next(ctx) {
        var user bson.M
        if err := cur.Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        delete(user, "password")
        users = append(users, user)
    }
     res_result  := map[string]interface{}{
        "ok": true,
        "status": "success",
        "users": users,
    }

    if err := cur.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set response header
    w.Header().Set("Content-Type", "application/json")

    // Encoding all users to JSON and sending the response
    if err := json.NewEncoder(w).Encode(res_result); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Println("Endpoint Hit: GetUsers")
}