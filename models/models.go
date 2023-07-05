package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	Role1 = "admin"
	Role2 = "subadmin"
	Role3 = "user"
)

type Users struct {
	Id        uuid.UUID
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedBy uuid.UUID
}

type UserRoles struct {
	Id   uuid.UUID
	Role string `json:"role"`
}

type Address struct {
	Id        uuid.UUID
	Name      string  `json:"addressName"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    uuid.UUID
}

type Restaurant struct {
	Id        uuid.UUID
	Name      string  `json:"restaurantName"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedBy uuid.UUID
}

type Dishes struct {
	Name         string `json:"dishName"`
	Price        int    `json:"price"`
	RestaurantId uuid.UUID
	CreatedBy    uuid.UUID
}

type Claims struct {
	UserID uuid.UUID `json:"userID"`
	jwt.StandardClaims
}
