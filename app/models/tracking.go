package models

import "time"

type Tracking struct {
    Sno              int       `gorm:"primaryKey;autoIncrement"`
    TrackingStatus   string    // Add this field if needed
    RecipientId        int
    CurrentLocation int
    EstimatedDelivery time.Time `gorm:"not null"`  // Add this field
    CreatedDt        time.Time
    UpdatedDt        time.Time
}