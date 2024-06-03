package routes

import (
	"hospital_management_portal/controllers"
	middleware "hospital_management_portal/middlewares"

	"net/http"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func DoctorRoutes(router *mux.Router) {

	//Doctor-Patient Routes 
	createPatientRoute := middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.CreatePatient), []string{"Doctor"})
    router.Handle("/doctor/patients", createPatientRoute).Methods("POST")

	updatePatientRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.UpdatePatient), []string{"Doctor"})
    router.Handle("/doctor/patients/{patientId}", updatePatientRoute).Methods("PUT")

	getPatientsRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetPatients), []string{"Doctor"})
    router.Handle("/doctor/patients", getPatientsRoute).Methods("GET")

	deletePatientRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.DeletePatient), []string{"Doctor"})
    router.Handle("/doctor/patients/{patientId}", deletePatientRoute).Methods("DELETE")

	getAllPatientAppointmentRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetAllAppointments), []string{"Doctor"})
    router.Handle("/doctor/appointments", getAllPatientAppointmentRoute).Methods("GET")



}
