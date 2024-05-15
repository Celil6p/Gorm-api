package handlers

import (
    "encoding/json"
    "gormTest/models"
    "net/http"
    "strconv"

    "github.com/cloudinary/cloudinary-go"
    "github.com/cloudinary/cloudinary-go/api/uploader"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

type PageHandler struct {
    db         *gorm.DB
    cloudinary *cloudinary.Cloudinary
}

func NewPageHandler(db *gorm.DB, cloudinaryConfig cloudinary.Config) *PageHandler {
    cld, _ := cloudinary.NewFromParams(cloudinaryConfig.CloudName, cloudinaryConfig.APIKey, cloudinaryConfig.APISecret)
    return &PageHandler{db: db, cloudinary: cld}
}

func (h *PageHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
    var page models.Page
    err := json.NewDecoder(r.Body).Decode(&page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get the magazineID from the URL parameters
    params := mux.Vars(r)
    magazineID, err := strconv.Atoi(params["magazineID"])
    if err != nil {
        http.Error(w, "Invalid magazine ID", http.StatusBadRequest)
        return
    }

    // Upload the image to Cloudinary
    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    uploadResult, err := h.cloudinary.Upload.Upload(r.Context(), file, uploader.UploadParams{
        PublicID: header.Filename,
        Folder:   "magazine_pages",
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    page.ImageURL = uploadResult.SecureURL
    page.MagazineID = uint(magazineID)

    result := h.db.Create(&page)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(page)
}

func (h *PageHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
    // ... similar implementation to CreatePage ...
}

func (h *PageHandler) DeletePage(w http.ResponseWriter, r *http.Request) {
    // ... delete the page from the database ...

    // Delete the image from Cloudinary
    publicID := "magazine_pages/" + header.Filename
    _, err = h.cloudinary.Upload.Destroy(r.Context(), uploader.DestroyParams{
        PublicID: publicID,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}