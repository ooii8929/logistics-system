package models

import "time"

type Recipient struct {
    ID       int       `gorm:"primaryKey;autoIncrement"`
    Name     string
    Address  string
    Phone    string
    CreatedDt  time.Time `gorm:"column:created_dt"`
}