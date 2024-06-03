package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"hospital_management_portal/database"
	"hospital_management_portal/models"
	request "hospital_management_portal/requests"
	"hospital_management_portal/response"
	"log"
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

func CreatePatient(w http.ResponseWriter, r *http.Request) {

	 //let's first validate our request
	 var patientCreatReq request.PatientRequest
	 decoder := json.NewDecoder(r.Body)
	 if err := decoder.Decode(&patientCreatReq); err != nil {
		 http.Error(w, "Invalid input", http.StatusBadRequest)
		 return
	 }
 
	 // Initialize the validator and validate the request data
	 validate := validator.New()
	 err := validate.Struct(patientCreatReq)
	 if err != nil {
		 if errs, ok := err.(validator.ValidationErrors); ok {
			 errMessages := request.CustomeErrorMessage(errs)
			 http.Error(w, strings.Join(errMessages, ", "), http.StatusBadRequest)
			 return
		 }
		 http.Error(w, "Validation failed", http.StatusInternalServerError)
		 return
	 }
	 
	 newPatient := models.Patient{
		ID: uuid.NewString(), 
		FirstName: patientCreatReq.FirstName,
		LastName: patientCreatReq.LastName,
		Phone: patientCreatReq.Phone,
		Email: patientCreatReq.Email,
        CreatedAt: time.Now(), 
        UpdatedAt: time.Now(),
    }

	  //lets first open the collection we need to insert the user in
	  collection:= database.OpenCollection(database.Client,"patients")
	  ctx,cancel :=context.WithTimeout(context.Background(),10*time.Second)
	  defer cancel()
	  //let's now insert the user
	  _,err=collection.InsertOne(ctx,newPatient)
	  if err != nil {
  
		  err_response := response.ErrorResponse{
			  OK: false,
			  Message: err.Error(),
			  Status: "failed",
		  }
		  json.NewEncoder(w).Encode(err_response)
	  }

		w.Header().Set("Content-Type", "application/json")
	
		resResult := map[string]interface{}{
			"ok":     true,
			"status": "success",
			"message": fmt.Sprintf(".%s %s added successfully", newPatient.FirstName,newPatient.LastName),
		}
		json.NewEncoder(w).Encode(resResult)

}

func UpdatePatient(w http.ResponseWriter, r *http.Request) {
    // Check and validate the request body
    var patientUpdateReq request.PatientRequest
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&patientUpdateReq); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    
    // Initialize the validator and validate the request data
    validate := validator.New()
    if err := validate.Struct(patientUpdateReq); err != nil {
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
    collection := database.OpenCollection(database.Client, "patients")
    vars := mux.Vars(r)
    patientID := vars["patientId"]

    // Create a filter to find the doctor by ID
    filter := bson.M{"id": patientID}

    // Check if the doctor exists
    var patient bson.M
    if err := collection.FindOne(ctx, filter).Decode(&patient); err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No Patient found with given ID", http.StatusNotFound)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    update := bson.M{
        "$set": bson.M{
            "first_name": patientUpdateReq.FirstName,
            "last_name": patientUpdateReq.LastName,
            "email": patientUpdateReq.Email,
            "phone": patientUpdateReq.Phone,
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
    json.NewEncoder(w).Encode(bson.M{"ok": true, "message": fmt.Sprintf("%s %s updated successfully", patientUpdateReq.FirstName, patientUpdateReq.LastName),"status":"success"})
}

func GetPatients(w http.ResponseWriter, r *http.Request) {
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

    collection := database.OpenCollection(database.Client, "patients")
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

    var patients []bson.M
    for cur.Next(ctx) {
        var patient bson.M
        if err := cur.Decode(&patient); err != nil {
            http.Error(w, "Failed to decode document", http.StatusInternalServerError)
            return
        }
        patients = append(patients, patient)
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
        "patients":      patients,
        "pagination": pagination, // Include the populated pagination map here
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res_result); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}


func DeletePatient(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Access the MongoDB collection
    collection := database.OpenCollection(database.Client, "patients")
    vars := mux.Vars(r)
    patientID := vars["patientId"]

    // Create a filter to find and delete the doctor by ID
    filter := bson.M{"id": patientID}

    // Perform the delete operation
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if the document was actually deleted
    if result.DeletedCount == 0 {
        http.Error(w, "No patient found with given ID", http.StatusNotFound)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{"ok": true, "message": "Patient deleted successfully", "status": "success"})
}

func GetAvailableDoctors(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := database.OpenCollection(database.Client, "doctors")
    filter := bson.M{"available": true}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        http.Error(w, "Failed to fetch doctors", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var doctors [] models.Doctor
    if err = cursor.All(ctx, &doctors); err != nil {
        http.Error(w, "Failed to parse doctors", http.StatusInternalServerError)
        return
    }
    res_result := map[string]interface{}{
        "ok":         true,
        "status":     "success",
        "doctors":      doctors,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res_result)
}


func ScheduleAppointment(w http.ResponseWriter, r *http.Request) {
    // Decode the request body into an appointment request struct
    var appointmentReq request.AppointmentRequest
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields() // This prevents accepting requests with unknown fields
    if err := decoder.Decode(&appointmentReq); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Initialize the validator and validate the request data
    validate := validator.New()
    if err := validate.Struct(appointmentReq); err != nil {
        if errs, ok := err.(validator.ValidationErrors); ok {
            errMessages := request.CustomeErrorMessage(errs)
            http.Error(w, strings.Join(errMessages, ", "), http.StatusBadRequest)
            return
        }
        http.Error(w, "Validation failed", http.StatusInternalServerError)
        return
    }

    // Set a timeout for database operations
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Open a collection and insert the new appointment
    collection := database.OpenCollection(database.Client, "appointments")
    _, err := collection.InsertOne(ctx, models.Appointment{
        ID:         uuid.NewString(), // Generate a new UUID for the appointment ID
        DoctorID:   appointmentReq.DoctorID,
        PatientID:  appointmentReq.PatientID,
        Time:       appointmentReq.Time,
        Completed:  false, // Initialize as false; update when the appointment is completed
    })
    if err != nil {
        http.Error(w, "Failed to schedule appointment", http.StatusInternalServerError)
        return
    }

    // Successful response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{"ok": true, "message": "Appointment scheduled successfully"})
}

func GetPatientAppointments(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    vars := mux.Vars(r)
    patientID := vars["patientId"]  // Directly use the UUID string

    // Validate UUID format
    if _, err := uuid.Parse(patientID); err != nil {
        http.Error(w, "Invalid patient ID format", http.StatusBadRequest)
        return
    }

    collection := database.OpenCollection(database.Client, "appointments")
    filter := bson.M{"patient_id": patientID}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var appointments []models.Appointment
    if err = cursor.All(ctx, &appointments); err != nil {
        http.Error(w, "Failed to parse appointments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{"ok": true, "status": "success","appointment":appointments})
}

func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Open the appointments collection
    collection := database.OpenCollection(database.Client, "appointments")

    // Find all appointments without filtering by patient_id
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Printf("Error fetching all appointments: %v", err)
        http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var appointments []models.Appointment
    if err := cursor.All(ctx, &appointments); err != nil {
        log.Printf("Error parsing appointments: %v", err)
        http.Error(w, "Failed to parse appointments", http.StatusInternalServerError)
        return
    }

    // Respond with all appointments
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bson.M{
        "ok": true,
        "status": "success",
        "appointments": appointments,
    })
}
