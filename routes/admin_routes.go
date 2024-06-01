package routes

import (
	"hospital_management_portal/controllers"
	middleware "hospital_management_portal/middlewares"

	"net/http"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func AdminRoutes(router *mux.Router) {
    createDoctorRoute := middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.CreateDoctor), []string{"Admin"})
    router.Handle("/admin/doctors", createDoctorRoute).Methods("POST")

	updateDoctorRoute:=middleware.RoleBasedJWTMiddleware(http.HandlerFunc(controllers.UpdateDoctor), []string{"Admin"})
    router.Handle("/admin/doctors/{doctorId}", updateDoctorRoute).Methods("PUT")


}