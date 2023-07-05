package server

import (
	"example/rms/handler"
	"example/rms/middlewares"
	"github.com/gorilla/mux"
)

func SetUpRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handler.LoginHandler).Methods("GET")

	r.Use(middlewares.JWTMiddleware)
	r.HandleFunc("/all-restaurants", handler.GetAllRestaurantsHandler).Methods("GET")
	r.HandleFunc("/all-dishes", handler.GetAllDishesHandler).Methods("GET")

	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(middlewares.MiddlewareUser)
	userRouter.HandleFunc("/distance", handler.CalcDistanceHandler).Methods("GET")
	userRouter.HandleFunc("/add-address", handler.InsertAddressHandler).Methods("POST")

	subadminRouter := r.PathPrefix("/subadmins").Subrouter()
	subadminRouter.Use(middlewares.MiddlewareSubAdmin)
	subadminRouter.HandleFunc("/all-restaurants-created", handler.GetAllRestaurantsBySubAdminHandler).Methods("GET")
	subadminRouter.HandleFunc("/all-users", handler.GetAllUsersBySubAdminsHandler).Methods("GET")
	subadminRouter.HandleFunc("/create-user", handler.CreateUserHandler).Methods("POST")
	subadminRouter.HandleFunc("/create-dish", handler.CreateDishHandler).Methods("POST")
	subadminRouter.HandleFunc("/create-restaurant", handler.CreateRestaurantHandler).Methods("POST")

	adminRouter := r.PathPrefix("/admins").Subrouter()
	adminRouter.Use(middlewares.MiddlewareAdmin)
	adminRouter.HandleFunc("/all-users", handler.GetAllUsersHandler).Methods("GET")
	adminRouter.HandleFunc("/create-user", handler.CreateUserHandler).Methods("POST")
	adminRouter.HandleFunc("/create-restaurant", handler.CreateRestaurantHandler).Methods("POST")
	adminRouter.HandleFunc("/create-subadmin", handler.CreateSubAdminHandler).Methods("POST")
	adminRouter.HandleFunc("/create-dish", handler.CreateDishHandler).Methods("POST")
	adminRouter.HandleFunc("/all-subadmins", handler.GetAllSubAdminsHandler).Methods("GET")

	return r
}
