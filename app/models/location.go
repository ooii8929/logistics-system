package models

type Location struct {
    LocationID int    `gorm:"primaryKey;autoIncrement"`
    Title      string
    City       string
    Address    string
}