package dbHelper

import (
	"example/rms/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetUserFromUsers(db *sqlx.DB, name, email string) (models.Users, error) {
	SQL := `SELECT id, password FROM users WHERE name = $1 AND email = $2`
	var user models.Users
	err := db.QueryRowx(SQL, name, email).Scan(&user.Id, &user.Password)
	return user, err
}

func CreateUser(db *sqlx.DB, user models.Users, createdby uuid.UUID) (uuid.UUID, error) {
	SQL := `INSERT INTO users(name, email, password) values ($1, $2, $3) RETURNING id`
	var userID uuid.UUID
	err := db.QueryRowx(SQL, user.Name, user.Email, user.Password).Scan(&userID)
	if err != nil || createdby == uuid.Nil {
		return userID, err
	}
	SQL = `UPDATE users SET created_by = $1 WHERE id = $2`
	_, err = db.Exec(SQL, createdby, userID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, err
}

func CreateUserRole(db *sqlx.DB, userID uuid.UUID, role string) error {
	SQL := `INSERT INTO user_roles (user_id, role_user) values ($1, $2)`
	_, err := db.Exec(SQL, userID, role)
	return err
}

func InsertAddress(db *sqlx.DB, address models.Address) error {
	SQL := `INSERT INTO address (name, latitude, longitude, user_id)
			VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(SQL, address.Name, address.Latitude, address.Longitude, address.UserID)
	return err
}

func GetUserRoles(db *sqlx.DB, userID uuid.UUID) ([]string, error) {
	SQL := `SELECT role_user FROM user_roles WHERE user_id = $1`
	var roles []string
	rows, err := db.Query(SQL, userID)
	for rows.Next() {
		var role string
		rows.Scan(&role)
		roles = append(roles, role)
	}
	return roles, err
}

func CreateRestaurant(db *sqlx.DB, restaurant models.Restaurant) error {
	SQL := `INSERT INTO restaurants (name, latitude, longitude, created_by)
            VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(SQL, restaurant.Name, restaurant.Latitude, restaurant.Longitude, restaurant.CreatedBy)
	return err
}

func CreateDish(db *sqlx.DB, dish models.Dishes) error {
	SQL := `INSERT INTO dishes (name, price, restaurant_id, created_by) 
			VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(SQL, dish.Name, dish.Price, dish.RestaurantId, dish.CreatedBy)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]models.Users, error) {
	SQL := `SELECT a.id, a.name, a.email, a.created_by 
			FROM users a INNER JOIN user_roles b
			ON a.id = b.user_id
			WHERE b.role_user = $1`
	var users []models.Users
	rows, err := db.Query(SQL, models.Role3)
	if err != nil {
		return []models.Users{}, err
	}
	for rows.Next() {
		var user models.Users
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedBy)
		users = append(users, user)
	}
	return users, nil
}

func GetAllUsersBySubadmin(db *sqlx.DB, userID uuid.UUID) ([]models.Users, error) {
	SQL := `SELECT id, name, email, created_by
			FROM users 
			WHERE created_by = $1`
	var users []models.Users
	rows, err := db.Query(SQL, userID)
	if err != nil {
		return []models.Users{}, err
	}
	for rows.Next() {
		var user models.Users
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedBy)
		users = append(users, user)
	}
	return users, nil
}

func GetAllSubAdmins(db *sqlx.DB) ([]models.Users, error) {
	SQL := `SELECT a.id, a.name, a.email
			FROM users a INNER JOIN user_roles b 
			ON a.id = b.user_id
			WHERE b.role_user = $1`
	rows, err := db.Query(SQL, models.Role2)
	if err != nil {
		return nil, err
	}
	var users []models.Users
	for rows.Next() {
		var user models.Users
		rows.Scan(&user.Id, &user.Name, &user.Email)
		users = append(users, user)
	}
	return users, nil
}

func GetAllRestaurants(db *sqlx.DB) ([]models.Restaurant, error) {
	SQL := `SELECT name, latitude, longitude 
			FROM restaurants`
	rows, err := db.Query(SQL)
	if err != nil {
		return []models.Restaurant{}, err
	} //todo: use rows.close
	defer rows.Close()
	var restaurants []models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		rows.Scan(&restaurant.Name, &restaurant.Latitude, &restaurant.Longitude)
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, nil
}

func GetAllRestaurantsBySubadmin(db *sqlx.DB, userID uuid.UUID) ([]models.Restaurant, error) {
	SQL := `SELECT name, created_by 
			FROM restaurants 
			WHERE created_by = $1`
	rows, err := db.Query(SQL, userID)
	if err != nil {
		return []models.Restaurant{}, err
	}
	var restaurants []models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		rows.Scan(&restaurant.Name, &restaurant.CreatedBy)
		restaurants = append(restaurants, restaurant)
	}
	return restaurants, nil
}

func GetRestaurantID(db *sqlx.DB, name string) (uuid.UUID, error) {
	SQL := `SELECT id 
			FROM restaurants 
			WHERE name = $1`
	var restaurantID uuid.UUID
	err := db.QueryRowx(SQL, name).Scan(&restaurantID)
	if err != nil {
		return uuid.Nil, err
	}
	return restaurantID, nil
}

func GetAllDishes(db *sqlx.DB, restaurantID uuid.UUID) ([]models.Dishes, error) {
	SQL := `SELECT name, price 
			FROM dishes
			WHERE restaurant_id = $1`
	rows, err := db.Query(SQL, restaurantID)
	if err != nil {
		return nil, err
	}
	var dishes []models.Dishes
	for rows.Next() {
		var dish models.Dishes
		rows.Scan(&dish.Name, &dish.Price)
		dishes = append(dishes, dish)
	}
	return dishes, nil
}

func GetAddressId(db *sqlx.DB, userID uuid.UUID, addressName string) (uuid.UUID, error) {
	SQL := `SELECT id FROM address
			WHERE user_id = $1 AND name = $2`
	var addressID uuid.UUID
	err := db.QueryRowx(SQL, userID, addressName).Scan(&addressID)
	return addressID, err
}

func CalcDistance(db *sqlx.DB, addressID, restaurantID uuid.UUID) (int, error) {
	SQL := `SELECT (earth_distance(ll_to_earth(address.latitude, address.longitude),
                                  ll_to_earth(restaurants.latitude, restaurants.longitude)
                                  )/1609.344)::integer AS distance_miles
			FROM address, restaurants
			WHERE address.id = $1 AND restaurants.id = $2`
	var distance int
	err := db.QueryRowx(SQL, addressID, restaurantID).Scan(&distance)
	return distance, err
}
