package handlers

import (
    "encoding/json"
    "gormTest/models"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type AdminHandler struct {
    db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
    return &AdminHandler{db: db}
}

func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
    var admin models.Admin
    err := json.NewDecoder(r.Body).Decode(&admin)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Hash the password using bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }
    admin.Password = string(hashedPassword)

    result := h.db.Create(&admin)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Exclude the password from the response
    admin.Password = ""

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) GetAdmin(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    adminID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid admin ID", http.StatusBadRequest)
        return
    }

    // // Retrieve the adminID from the context
    // contextAdminID, ok := r.Context().Value("adminID").(int)
    // if !ok {
    //     http.Error(w, "Unauthorized", http.StatusUnauthorized)
    //     return
    // }

    // // Check if the requested adminID matches the authenticated adminID
    // if adminID != contextAdminID {
    //     http.Error(w, "Unauthorized", http.StatusUnauthorized)
    //     return
    // }

    var admin models.Admin
    if err := h.db.First(&admin, adminID).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    adminID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid admin ID", http.StatusBadRequest)
        return
    }

    // Retrieve the adminID from the context
    contextAdminID, ok := r.Context().Value("adminID").(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Check if the requested adminID matches the authenticated adminID
    if adminID != contextAdminID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var admin models.Admin
    result := h.db.First(&admin, adminID)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Admin not found", http.StatusNotFound)
        } else {
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }

    err = json.NewDecoder(r.Body).Decode(&admin)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result = h.db.Save(&admin)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(admin)
}

func (h *AdminHandler) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    adminID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid admin ID", http.StatusBadRequest)
        return
    }

    // Retrieve the adminID from the context
    contextAdminID, ok := r.Context().Value("adminID").(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Check if the requested adminID matches the authenticated adminID
    if adminID != contextAdminID {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    result := h.db.Delete(&models.Admin{}, adminID)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}