package main

import (
	"example/rms/database"
	"example/rms/server"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func main() {
	logrus.Println(os.Getenv("HOST"))
	if err := database.ConnectAndMigrate(
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("databaseName"),
		os.Getenv("user"),
		os.Getenv("password"),
		database.SSLModeDisable); err != nil {
		logrus.Fatalf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Infof("migration successful!!")

	r := server.SetUpRoutes()

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

//todo: remove business logic from handler layer
