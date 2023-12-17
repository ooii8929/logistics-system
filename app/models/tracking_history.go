package models

import "time"

type TrackingHistory struct {
    RecordID       int       `gorm:"primaryKey;autoIncrement"`
    TrackingStatus string
    Sno            int
    LocationID     int
    CreatedDt      time.Time
}
