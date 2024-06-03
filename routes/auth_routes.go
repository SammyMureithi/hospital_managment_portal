package routes

import (
	"hospital_management_portal/controllers"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func UserRoutes(router *mux.Router) {
  
     router.HandleFunc("/auth/register", controllers.SignUp).Methods("POST")
     router.HandleFunc("/auth/login", controllers.Login).Methods("POST")
}
