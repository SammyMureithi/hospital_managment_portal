package request

import (
	"time"
)

type AppointmentRequest struct {
	

	DoctorID  string        `json:"doctor_id" bson:"doctor_id" validate:"required"`
	PatientID string        ` json:"patient_id" bson:"patient_id" validate:"required"`
	Time      time.Time     `json:"time"bson:"time"`
	Completed bool          `bson:"completed"`
}