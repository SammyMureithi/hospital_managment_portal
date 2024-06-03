package routes

import (
	"hospital_management_portal/controllers"
	middleware "hospital_management_portal/middlewares"

	"net/http"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func PatientRoutes(router *mux.Router) {

	//Doctor-Patient Routes 
	createPatientAppointmentRoute := middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.ScheduleAppointment), []string{"Patient"})
    router.Handle("/patient/appointments", createPatientAppointmentRoute).Methods("POST")

	viewAvailableDoctors:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetAvailableDoctors), []string{"Patient"})
    router.Handle("/patient/doctors", viewAvailableDoctors).Methods("GET")

	getPatientsAppointmentsRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetPatientAppointments), []string{"Patient"})
    router.Handle("/patient/appointments/{patientId}", getPatientsAppointmentsRoute).Methods("GET")


}
