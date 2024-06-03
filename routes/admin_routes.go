package routes

import (
	"hospital_management_portal/controllers"
	middleware "hospital_management_portal/middlewares"

	"net/http"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func AdminRoutes(router *mux.Router) {
	//Admin-Doctors Routes
    createDoctorRoute := middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.CreateDoctor), []string{"Admin"})
    router.Handle("/admin/doctors", createDoctorRoute).Methods("POST")

	updateDoctorRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.UpdateDoctor), []string{"Admin"})
    router.Handle("/admin/doctors/{doctorId}", updateDoctorRoute).Methods("PUT")

	getDoctorsRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetDoctors), []string{"Admin"})
    router.Handle("/admin/doctors", getDoctorsRoute).Methods("GET")

	deleteDoctorsRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.DeleteDoctor), []string{"Admin"})
    router.Handle("/admin/doctors/{doctorId}", deleteDoctorsRoute).Methods("DELETE")


	//Admin-Patient Routes 
	createPatientRoute := middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.CreatePatient), []string{"Admin"})
    router.Handle("/admin/patients", createPatientRoute).Methods("POST")

	updatePatientRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.UpdatePatient), []string{"Admin"})
    router.Handle("/admin/patients/{patientId}", updatePatientRoute).Methods("PUT")

	getPatientsRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetPatients), []string{"Admin"})
    router.Handle("/admin/patients", getPatientsRoute).Methods("GET")

	deletePatientRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.DeletePatient), []string{"Admin"})
    router.Handle("/admin/patients/{patientId}", deletePatientRoute).Methods("DELETE")

	getAllPatientAppointmentRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.GetAllAppointments), []string{"Doctor"})
    router.Handle("/admin/appointments", getAllPatientAppointmentRoute).Methods("GET")



}
