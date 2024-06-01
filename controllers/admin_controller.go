package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"hospital_management_portal/database"
	"hospital_management_portal/models"
	request "hospital_management_portal/requests"
	"hospital_management_portal/response"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
        CreatedAt: time.Now(), 
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

func UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    // Check and validate the request body
    var docUpdateReq request.DoctorRequest
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&docUpdateReq); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    
    // Initialize the validator and validate the request data
    validate := validator.New()
    if err := validate.Struct(docUpdateReq); err != nil {
        if errs, ok := err.(validator.ValidationErrors); ok {
            errMessages := request.CustomeErrorMessage(errs)
            http.Error(w, strings.Join(errMessages, ", "), http.StatusBadRequest)
            return
        }
        http.Error(w, "Validation failed", http.StatusInternalServerError)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Access the MongoDB collection
    collection := database.OpenCollection(database.Client, "doctors")
    vars := mux.Vars(r)
    doctorID := vars["doctorId"]

    // Create a filter to find the doctor by ID
    filter := bson.M{"id": doctorID}

    // Check if the doctor exists
    var doctor bson.M
    if err := collection.FindOne(ctx, filter).Decode(&doctor); err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No doctor found with given ID", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    update := bson.M{
        "$set": bson.M{
            "first_name": docUpdateReq.FirstName,
            "last_name": docUpdateReq.LastName,
            "email": docUpdateReq.Email,
            "phone": docUpdateReq.Phone,
            "updated_at": time.Now(),
        },
    }

    // Update the doctor document
    if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
        http.Error(w, "Failed to update doctor", http.StatusInternalServerError)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{"ok": true, "message": fmt.Sprintf("Dr. %s %s updated successfully", docUpdateReq.FirstName, docUpdateReq.LastName),"status":"success"})
}

func GetDoctors(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    query := r.URL.Query()
    limitQuery := query.Get("limit")
    pageQuery := query.Get("page")

    limit := 10
    page := 1

    if l, err := strconv.Atoi(limitQuery); err == nil && l > 0 {
        limit = l
    }
    if p, err := strconv.Atoi(pageQuery); err == nil && p > 1 {
        page = p
    }

    skip := (page - 1) * limit

    collection := database.OpenCollection(database.Client, "doctors")
    total, err := collection.CountDocuments(ctx, bson.D{})
    if err != nil {
        http.Error(w, "Failed to count documents", http.StatusInternalServerError)
        return
    }

    totalPages := int(math.Ceil(float64(total) / float64(limit)))

    findOptions := options.Find()
    findOptions.SetLimit(int64(limit))
    findOptions.SetSkip(int64(skip))

    cur, err := collection.Find(ctx, bson.D{}, findOptions)
    if err != nil {
        http.Error(w, "Failed to find documents", http.StatusInternalServerError)
        return
    }
    defer cur.Close(ctx)

    var doctors []bson.M
    for cur.Next(ctx) {
        var doctor bson.M
        if err := cur.Decode(&doctor); err != nil {
            http.Error(w, "Failed to decode document", http.StatusInternalServerError)
            return
        }
        doctors = append(doctors, doctor)
    }

    // Create and populate the pagination map
    pagination := map[string]interface{}{
        "current_page": page,
        "total_pages":  totalPages,
        "limit":        limit,
        "total_items":  total,
    }

    // Add URLs to the pagination map before adding it to the result
    if page < totalPages {
        nextURL := fmt.Sprintf("%s?limit=%d&page=%d", r.URL.Path, limit, page+1)
        pagination["next_page_url"] = nextURL
    }
    if page > 1 {
        prevURL := fmt.Sprintf("%s?limit=%d&page=%d", r.URL.Path, limit, page-1)
        pagination["previous_page_url"] = prevURL
    }

    // Create the final response map
    res_result := map[string]interface{}{
        "ok":         true,
        "status":     "success",
        "users":      doctors,
        "pagination": pagination, // Include the populated pagination map here
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res_result); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}


func DeleteDoctor(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Access the MongoDB collection
    collection := database.OpenCollection(database.Client, "doctors")
    vars := mux.Vars(r)
    doctorID := vars["doctorId"]

    // Create a filter to find and delete the doctor by ID
    filter := bson.M{"id": doctorID}

    // Perform the delete operation
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if the document was actually deleted
    if result.DeletedCount == 0 {
        http.Error(w, "No doctor found with given ID", http.StatusNotFound)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{"ok": true, "message": "Doctor deleted successfully", "status": "success"})
}


