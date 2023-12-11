package main

import (
    "log"
    "net/http"
    "time"
    "fmt"

    // Fake API
    "math/rand"
    "strconv"

    // Query API
    "encoding/json"

    // DB struct
    "logistics-track/models"
    "logistics-track/helpers"
    "logistics-track/database"
    "logistics-track/report"

    // Redis connection
    "logistics-track/redis"

    "gorm.io/gorm"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gin-gonic/gin"
)

type TrackingStatus string

const (
    Created           TrackingStatus = "Created"
    PackageReceived   TrackingStatus = "Package received"
    InTransit         TrackingStatus = "In transit"
    OutForDelivery    TrackingStatus = "Out for delivery"
    DeliveryAttempted TrackingStatus = "Delivery attempted"
    Delivered         TrackingStatus = "Delivered"
    ReturnedToSender  TrackingStatus = "Returned to sender"
    Exception         TrackingStatus = "Exception"
)

// Add your other type definitions and createTables functions here
type ApiResponse struct {
    Status string      `json:"status"`
    Error  interface{} `json:"error"`
}

// Add this type definition
type ResponseData struct {
    Status          string             `json:"status"`
    Data            *models.Tracking          `json:"data,omitempty"`
    Details         []models.TrackingHistory  `json:"details,omitempty"`
    Recipient       *models.Recipient         `json:"recipient,omitempty"`
    CurrentLocation *models.Location          `json:"current_location,omitempty"`
    Error           string             `json:"error,omitempty"`
}

func randomString(n int) string {
    // Update this line
    return helpers.RandomString(n)
}


func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func queryAPI(db *gorm.DB, sno int) ResponseData {
    var tracking models.Tracking
    if err := db.Where("sno = ?", sno).First(&tracking).Error; err != nil {
        return ResponseData{
            Status: "error",
            Error:  fmt.Sprintf("No tracking information found for sno: %d", sno),
        }
    }

    var recipient models.Recipient
    db.First(&recipient, tracking.RecipientId)

    var currentLocation models.Location
    db.Where("location_id = ?", tracking.CurrentLocation).First(&currentLocation)

    var details []models.TrackingHistory
    db.Where("sno = ?", sno).Find(&details)

    return ResponseData{
        Status:          "success",
        Data:            &tracking,
        Details:         details,
        Recipient:       &recipient,
        CurrentLocation: &currentLocation,
        Error:           "",
    }
}


func GenerateFakeData(c *gin.Context, num int) {
    db := database.CreateDatabaseAndServe()

    allStatuses := []TrackingStatus{Created, PackageReceived, InTransit, OutForDelivery, DeliveryAttempted, Delivered, ReturnedToSender, Exception}
    trackingList := make([]models.Tracking, 0)

    if num > 0 {
        // Insert fake locations
        for i := 0; i < num; i++ {
            loc := models.Location{Title: randomString(10), City: randomString(5), Address: randomString(15)}
            db.Create(&loc)
        }

        // Insert fake recipients
        for i := 0; i < num; i++ {
            rec := models.Recipient{Name: randomString(5), Address: randomString(15), Phone: randomString(10), CreatedDt: time.Now()}
            db.Create(&rec)
        }

        // Insert fake trackings
        for i := 1; i < num; i++ {
            tr := models.Tracking{
                TrackingStatus:    string(allStatuses[rand.Intn(len(allStatuses))]),
                RecipientId:         rand.Intn(num) + 1,
                CurrentLocation:   rand.Intn(num)+1,
                EstimatedDelivery: time.Now().Add(time.Duration(rand.Intn(168)) * time.Hour),
                CreatedDt:         time.Now(),
                UpdatedDt:         time.Now(),
            }
            db.Create(&tr)
            trackingList = append(trackingList, tr)
        }

        // Insert fake tracking histories
        for i := 0; i < num; i++ {
            trhist := models.TrackingHistory{
                TrackingStatus: string(allStatuses[rand.Intn(len(allStatuses))]),
                Sno:            rand.Intn(num) + 1,
                LocationID:     rand.Intn(num) + 1,
                CreatedDt:      time.Now(),
            }
            db.Create(&trhist)
        }
    }
    c.JSON(http.StatusOK, trackingList)
}

func main() {

    redisAddr := "redis-server:6379"
    redisPassword := ""
    redisDB := 0

    client, err := redis.ConnectRedis(redisAddr, redisPassword, redisDB)
    if err != nil {
        panic(err)
    }

    defer client.Close()

    r := gin.Default()

    r.GET("/hello", func(c *gin.Context) {
        response := gin.H{
            "message": "Hello, World!",
        }
        c.JSON(http.StatusOK, response)
    })

    r.GET("/fake", func(c *gin.Context) {
        numStr := c.DefaultQuery("num", "0")
        num, _ := strconv.Atoi(numStr)
        
        if num > 0 {
            GenerateFakeData(c, num)
        } else {
            response := ApiResponse{Status: "error", Error: nil}
            c.JSON(http.StatusBadRequest, response)
        }
    })

    r.GET("/generate_report", func(c *gin.Context) {
        summary, err := report.GetTrackingSummary()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "error":  fmt.Sprintf("Failed to get tracking summary: %v", err),
            })
            return
        }
    
        // 使用查询结果构建响应
        response := gin.H{
            "created_dt":        time.Now().Format(time.RFC3339),
            "trackingSummary":   summary,
        }
        c.JSON(http.StatusOK, response)
    })

    r.GET("/query", func(c *gin.Context) {
        snoStr := c.Query("sno")
        snoInt, err := strconv.Atoi(snoStr)
        if err != nil {
            response := ApiResponse{Status: "error", Error: "Invalid 'sno' query parameter"}
            c.JSON(http.StatusBadRequest, response)
            return
        }

        db := database.CreateDatabaseAndServe()
        responseData := queryAPI(db, snoInt)

        if responseData.Error != "" {
            c.JSON(http.StatusNotFound, responseData)
        } else {
            c.JSON(http.StatusOK, responseData)
        }
    })

    if err := r.Run(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}