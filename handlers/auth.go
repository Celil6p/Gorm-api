package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "gormTest/models"
    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type AuthHandler struct {
    db        *gorm.DB
    jwtSecret []byte
}

func NewAuthHandler(db *gorm.DB, jwtSecret []byte) *AuthHandler {
    return &AuthHandler{db: db, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var admin models.Admin
    if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Retrieve the admin from the database based on the email
    var dbAdmin models.Admin
    if err := h.db.Where("email = ?", admin.Email).First(&dbAdmin).Error; err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Compare the provided password with the hashed password in the database
    if err := bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)); err != nil {
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Generate a JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "adminID": dbAdmin.ID,
        "exp":     time.Now().Add(30 * 24 * time.Hour).Unix(), // Token expires in 30 days
    })

    tokenString, err := token.SignedString(h.jwtSecret)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return the token as the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}