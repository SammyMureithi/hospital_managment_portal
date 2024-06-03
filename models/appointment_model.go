package models

import (
	"time"
)

type Appointment struct {
	ID        string `bson:"_id,omitempty"`
	DoctorID  string `bson:"doctor_id"`
	PatientID string `bson:"patient_id"`
	Time      time.Time          `bson:"time"`
	Completed bool               `bson:"completed"`
}