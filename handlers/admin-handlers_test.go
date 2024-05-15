package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "gormTest/models"

    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    // Set up a test database connection
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Europe/Istanbul"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to the test database")
    }
    return db
}

func cleanupTestDB(db *gorm.DB) {
    db.Exec("DELETE FROM admins")
    db.Exec("ALTER SEQUENCE admins_id_seq RESTART WITH 1")
}

func TestCreateAdmin(t *testing.T) {
    db := setupTestDB()
    defer cleanupTestDB(db)
    adminHandler := NewAdminHandler(db)

    admin := models.Admin{
        Email:    "test@example.com",
        Password: "password",
    }

    body, _ := json.Marshal(admin)
    req, _ := http.NewRequest("POST", "/admins", bytes.NewBuffer(body))
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(adminHandler.CreateAdmin)
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code)

    var createdAdmin models.Admin
    err := json.Unmarshal(rr.Body.Bytes(), &createdAdmin)
    assert.NoError(t, err)
    assert.Equal(t, admin.Email, createdAdmin.Email)
    assert.NotEmpty(t, createdAdmin.ID)
}

func TestGetAdmin(t *testing.T) {
    db := setupTestDB()
    defer cleanupTestDB(db)
    adminHandler := NewAdminHandler(db)

    // Create a test admin in the database
    admin := models.Admin{
        Email:    "test@example.com",
        Password: "password",
    }
    db.Create(&admin)

    req, _ := http.NewRequest("GET", "/admins/{id}", nil)
    req = mux.SetURLVars(req, map[string]string{
        "id": "1",
    })
    req = req.WithContext(context.WithValue(req.Context(), "admin_id", int(admin.ID)))

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(adminHandler.GetAdmin)
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var retrievedAdmin models.Admin
    err := json.Unmarshal(rr.Body.Bytes(), &retrievedAdmin)
    assert.NoError(t, err)
    assert.Equal(t, admin.ID, retrievedAdmin.ID)
    assert.Equal(t, admin.Email, retrievedAdmin.Email)
}

func TestUpdateAdmin(t *testing.T) {
    db := setupTestDB()
    defer cleanupTestDB(db)
    adminHandler := NewAdminHandler(db)

    // Create a test admin in the database
    admin := models.Admin{
        Email:    "test@example.com",
        Password: "password",
    }
    db.Create(&admin)

    updatedAdmin := models.Admin{
        Email: "updated@example.com",
    }
    body, _ := json.Marshal(updatedAdmin)

    req, _ := http.NewRequest("PUT", "/admins/{id}", bytes.NewBuffer(body))
    req = mux.SetURLVars(req, map[string]string{
        "id": "1",
    })
    req = req.WithContext(context.WithValue(req.Context(), "admin_id", int(admin.ID)))

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(adminHandler.UpdateAdmin)
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var updatedAdminResponse models.Admin
    err := json.Unmarshal(rr.Body.Bytes(), &updatedAdminResponse)
    assert.NoError(t, err)
    assert.Equal(t, updatedAdmin.Email, updatedAdminResponse.Email)
}

func TestDeleteAdmin(t *testing.T) {
    db := setupTestDB()
    defer cleanupTestDB(db)
    adminHandler := NewAdminHandler(db)

    // Create a test admin in the database
    admin := models.Admin{
        Email:    "test@example.com",
        Password: "password",
    }
    db.Create(&admin)

    req, _ := http.NewRequest("DELETE", "/admins/{id}", nil)
    req = mux.SetURLVars(req, map[string]string{
        "id": "1",
    })
    req = req.WithContext(context.WithValue(req.Context(), "admin_id", int(admin.ID)))

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(adminHandler.DeleteAdmin)
    handler.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNoContent, rr.Code)

    var deletedAdmin models.Admin
    result := db.First(&deletedAdmin, admin.ID)
    assert.Error(t, result.Error)
    assert.EqualError(t, result.Error, gorm.ErrRecordNotFound.Error())
}