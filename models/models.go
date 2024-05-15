package models

import (
    "gorm.io/gorm"
)

type Admin struct {
    gorm.Model
    Email     string     `gorm:"unique"`
    Password  string
    Magazines []Magazine
}

type Magazine struct {
    gorm.Model
    Title     string
    Pages     []Page
    AdminID   uint
    Admin     Admin
    Published bool
}

type Page struct {
    gorm.Model
    PageNumber int
    ImageURL   string
    MagazineID uint
    Magazine   Magazine
}