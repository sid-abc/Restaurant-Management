package handler

import (
	"database/sql"
	"encoding/json"
	"example/rms/database"
	"example/rms/database/dbHelper"
	"example/rms/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)
	tx, err := database.Todo.Beginx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := dbHelper.CreateUser(database.Todo, user, uuid.Nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	err = dbHelper.CreateUserRole(database.Todo, userID, models.Role3)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDatabase, err := dbHelper.GetUserFromUsers(database.Todo, user.Name, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDatabase.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &models.Claims{
		UserID: userDatabase.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"token": tokenString,
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func InsertAddressHandler(w http.ResponseWriter, r *http.Request) {
	var address models.Address
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	address.UserID = userID
	err = dbHelper.InsertAddress(database.Todo, address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)
	tx, err := database.Todo.Beginx()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id, err := dbHelper.CreateUser(database.Todo, user, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	err = dbHelper.CreateUserRole(database.Todo, id, models.Role3)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func CreateSubAdminHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	var user models.Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)
	id, err := dbHelper.CreateUser(database.Todo, user, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = dbHelper.CreateUserRole(database.Todo, id, models.Role2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := dbHelper.GetAllUsers(database.Todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetAllUsersBySubAdminsHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	users, err := dbHelper.GetAllUsersBySubadmin(database.Todo, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetAllSubAdminsHandler(w http.ResponseWriter, r *http.Request) {
	subadmins, err := dbHelper.GetAllSubAdmins(database.Todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subadmins)
}

func CreateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	var restaurant models.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	restaurant.CreatedBy = userID
	err = dbHelper.CreateRestaurant(database.Todo, restaurant)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func GetAllRestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	restaurants, err := dbHelper.GetAllRestaurants(database.Todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(restaurants)
}

func GetAllRestaurantsBySubAdminHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	restaurants, err := dbHelper.GetAllRestaurantsBySubadmin(database.Todo, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restaurants)
}

func CreateDishHandler(w http.ResponseWriter, r *http.Request) {
	var dish models.Dishes
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	restaurantName := r.URL.Query().Get("restaurantName")
	restaurantID, err := dbHelper.GetRestaurantID(database.Todo, restaurantName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	dish.RestaurantId = restaurantID
	dish.CreatedBy = userID
	err = dbHelper.CreateDish(database.Todo, dish)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetAllDishesHandler(w http.ResponseWriter, r *http.Request) {
	restaurantName := r.URL.Query().Get("restaurantName")
	restaurantID, err := dbHelper.GetRestaurantID(database.Todo, restaurantName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dishes, err := dbHelper.GetAllDishes(database.Todo, restaurantID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dishes)
}

func CalcDistanceHandler(w http.ResponseWriter, r *http.Request) {
	restaurantName := r.URL.Query().Get("restaurantName")
	address := r.URL.Query().Get("address")
	claims := r.Context().Value("claims").(*models.Claims)
	userID := claims.UserID
	restaurantID, err := dbHelper.GetRestaurantID(database.Todo, restaurantName)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	addressID, err := dbHelper.GetAddressId(database.Todo, userID, address)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	distance, err := dbHelper.CalcDistance(database.Todo, addressID, restaurantID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(distance)
}
