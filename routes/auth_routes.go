package routes

import (
	"hospital_management_portal/controllers"

	"github.com/gorilla/mux"
)

// UserRoutes function to initialize user routes
func UserRoutes(router *mux.Router) {
      // Apply middleware to the GET routes
    //   router.Handle("/users", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetUsers))).Methods("GET")
    //   router.Handle("/user", middleware.JWTMiddleware(http.HandlerFunc(controllers.GetUser))).Methods("GET")
      
     router.HandleFunc("/auth/register", controllers.SignUp).Methods("POST")
     router.HandleFunc("/auth/login", controllers.Login).Methods("POST")
}
