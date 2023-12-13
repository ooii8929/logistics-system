package main

import (
    "log"
    "net/http"
    "time"
    "fmt"

    // Fake API
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

// Add your other type definitions and createTables functions here
type ApiResponse struct {
    Status string      `json:"status"`
    Error  interface{} `json:"error"`
}

type ApiResponseWithTrackingList struct {
    Status   string             `json:"status"`
    Error    string             `json:"error,omitempty"`
    Trackings []models.Tracking `json:"trackings,omitempty"`
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

func main() {

    // Load environment variables.
    env := &Environment{}
    err := env.load()
    if err != nil {
        log.Printf("Error loading environment variables: %v", err)
        return
    }

    client, err := redis.ConnectRedis(env.RedisAddr, env.RedisPassword, env.RedisDB)
    if err != nil {
        panic(err)
    }

    defer client.Close()

    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        response := gin.H{
            "message": "Hello, World!",
        }
        c.JSON(http.StatusOK, response)
    })

    r.GET("/init-database", func(c *gin.Context) {

        err := database.CreateAndInitializeDatabase(env.DBUsername, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
        if err != nil {
            log.Printf("Failed to create and initialize database: %v", err)
    
            c.JSON(500, gin.H{
                "status": "failed",
                "error":  "Failed to create and initialize the database",
            })
            return
        }
    
        c.JSON(200, gin.H{
            "status": "success",
            "message": "Database successfully created and initialized",
        })
    })
    

    r.GET("/fake", func(c *gin.Context) {
        numStr := c.DefaultQuery("num", "0")
        num, _ := strconv.Atoi(numStr)

        db, err := database.ConnectToMySQLwithTable(env.DBUsername, env.DBPassword, env.DBHost, env.DBName, env.DBPort)
        if err != nil {
            log.Fatalf("Failed to connect to MySQL: %v", err)
        }

        if num > 0 {
            trackings := database.GenerateFakeData(db, num)
            response := ApiResponseWithTrackingList{
                Status:    "success",
                Trackings: trackings,
            }
            c.JSON(http.StatusOK, response)
        } else {
            response := ApiResponse{Status: "error", Error: "Invalid parameter"}
            c.JSON(http.StatusBadRequest, response)
        }
    })

    r.GET("/generate-report", func(c *gin.Context) {

        db, err := database.ConnectToMySQLwithTable(env.DBUsername, env.DBPassword, env.DBHost, env.DBName, env.DBPort)
        if err != nil {
            log.Fatalf("Failed to connect to MySQL: %v", err)
        }
        summary, err := report.GetTrackingSummary(db)
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

        db, err := database.ConnectToMySQLwithTable(env.DBUsername, env.DBPassword, env.DBHost, env.DBName, env.DBPort)
        if err != nil {
            log.Fatalf("Failed to connect to MySQL: %v", err)
        }
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