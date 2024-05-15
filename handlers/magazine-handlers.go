package handlers

import (
    "encoding/json"
    "gormTest/models"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

type MagazineHandler struct {
    db *gorm.DB
}

func NewMagazineHandler(db *gorm.DB) *MagazineHandler {
    return &MagazineHandler{db: db}
}

func (h *MagazineHandler) CreateMagazine(w http.ResponseWriter, r *http.Request) {
    var magazine models.Magazine
    err := json.NewDecoder(r.Body).Decode(&magazine)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get the authenticated admin ID from the context
    adminID, ok := r.Context().Value("adminID").(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    magazine.AdminID = uint(adminID)

    result := h.db.Create(&magazine)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(magazine)
}

func (h *MagazineHandler) GetMagazine(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    magazineID, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid magazine ID", http.StatusBadRequest)
        return
    }

    var magazine models.Magazine
    if err := h.db.Preload("Pages").First(&magazine, magazineID).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(magazine)
}

func (h *MagazineHandler) UpdateMagazine(w http.ResponseWriter, r *http.Request) {
    // ... implementation similar to UpdateAdmin ...
}

func (h *MagazineHandler) DeleteMagazine(w http.ResponseWriter, r *http.Request) {
    // ... implementation similar to DeleteAdmin ...
}