package models
import "time"

type Location struct {
    LocationID int    `gorm:"primaryKey;autoIncrement"`
    Title      string
    City       string
    Address    string
    CreatedDt  time.Time `gorm:"column:created_dt"`
}