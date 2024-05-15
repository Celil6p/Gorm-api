package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "gormTest/handlers"
    "gormTest/models"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gormTest/middlewares"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    jwtSecret := os.Getenv("JWT_SECRET")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul", dbHost, dbUser, dbPassword, dbName, dbPort)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    db.AutoMigrate(&models.Admin{}, &models.Magazine{}, &models.Page{})

    adminHandler := handlers.NewAdminHandler(db)
    authHandler := handlers.NewAuthHandler(db, []byte(jwtSecret))

    r := mux.NewRouter()
    r.HandleFunc("/admins", adminHandler.CreateAdmin).Methods("POST")
    r.HandleFunc("/admins/{id}", middlewares.AuthMiddleware([]byte(jwtSecret))(adminHandler.GetAdmin)).Methods("GET")
    r.HandleFunc("/admins/{id}", middlewares.AuthMiddleware([]byte(jwtSecret))(adminHandler.UpdateAdmin)).Methods("PUT")
    r.HandleFunc("/admins/{id}", middlewares.AuthMiddleware([]byte(jwtSecret))(adminHandler.DeleteAdmin)).Methods("DELETE")

    // Login route
    r.HandleFunc("/login", authHandler.Login).Methods("POST")

    log.Println("Server started on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}