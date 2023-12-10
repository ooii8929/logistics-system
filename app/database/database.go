package database

import (
    "log"
    "fmt"
    "time"
    
    "gorm.io/gorm"
    _ "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "logistics-track/models"
)

func CreateDatabaseAndServe() *gorm.DB {
    // 替換為您的 MySQL 用戶名、密碼和服務器地址
    username := "root"
    password := ""
    host := "my-mysql"
    port := 3306

    // 連接到 MySQL 服務器（不指定數據庫名稱）
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }

    // 創建新的數據庫
    databaseName := "logistics"
    err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", databaseName)).Error
    if err != nil {
        log.Fatalf("Failed to create database: %v", err)
    } else {
        fmt.Printf("Database '%s' created successfully.\n", databaseName)
    }

    // 斷開與未指定數據庫的 MySQL 連接
    sqlDB, _ := db.DB()
    sqlDB.Close()

    // 連接到創建的數據庫
    dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, databaseName)
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to MySQL with new database: %v", err)
    }

    CreateTables(db)
    db.AutoMigrate(&models.Location{}, &models.Recipient{}, &models.Tracking{}, &models.TrackingHistory{})
    
    // Insert data
    // Insert fake locations
    locations := []models.Location{
        {Title: "Warehouse 1", City: "New York", Address: "123 Main St"},
        {Title: "Warehouse 2", City: "San Francisco", Address: "456 Market St"},
    }

    for _, loc := range locations {
        db.Create(&loc)
    }

    // Insert fake recipients
    recipients := []models.Recipient{
        {Name: "Alice", Address: "789 Elm St", Phone: "555-1234", CreatedDt: time.Now()},
        {Name: "Bob", Address: "321 Oak St", Phone: "555-5678", CreatedDt: time.Now()},
    }
    
    for _, rec := range recipients {
        db.Create(&rec)
    }

    // Insert fake trackings
    trackings := []models.Tracking{
        {
            TrackingStatus:    "Active",
            RecipientId:         1,
            CurrentLocation: 1,
            EstimatedDelivery: time.Now().Add(2 * time.Hour), // 添加一个估计的交货时间
            CreatedDt:         time.Now(),
            UpdatedDt:         time.Now(),
        },
        {
            TrackingStatus:    "Completed",
            RecipientId:         2,
            CurrentLocation: 2,
            EstimatedDelivery: time.Now().Add(-2 * time.Hour), // 添加一个过去的交货时间
            CreatedDt:         time.Now(),
            UpdatedDt:         time.Now(),
        },
    }

    for _, tr := range trackings {
        db.Create(&tr)
    }

    // Insert fake tracking histories
    trackingHistories := []models.TrackingHistory{
        {TrackingStatus: "Picked up", Sno: 1, LocationID: 1, CreatedDt: time.Now()},
        {TrackingStatus: "In Transit", Sno: 2, LocationID: 2, CreatedDt: time.Now()},
    }

    for _, trhist := range trackingHistories {
        db.Create(&trhist)
    }

    return db
}

func CreateTables(db *gorm.DB) {
    tables := []string{
        `
        CREATE TABLE IF NOT EXISTS locations (
            location_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            title TEXT NOT NULL,
            city TEXT NOT NULL,
            address TEXT NOT NULL
        )`,
        `
        CREATE TABLE IF NOT EXISTS recipients (
            id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            name TEXT NOT NULL,
            address TEXT NOT NULL,
            phone TEXT NOT NULL,
            created_dt DATETIME NOT NULL
        )`,
        `
        CREATE TABLE IF NOT EXISTS trackings (
            sno INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            tracking_status TEXT NOT NULL,
            estimated_delivery DATETIME NOT NULL,
            recipient_id INT UNSIGNED,
            current_location INT UNSIGNED,
            created_dt DATETIME NOT NULL,
            updated_dt DATETIME NOT NULL
        )`,
        `
        CREATE TABLE IF NOT EXISTS tracking_histories (
            record_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            tracking_status TEXT NOT NULL,
            sno INT UNSIGNED,
            location_id INT UNSIGNED,
            created_dt DATETIME NOT NULL
        );`,
    }

    for _, table := range tables {
        err := db.Exec(table).Error
        if err != nil {
            log.Fatalf("Failed to create table: %v", err)
        }
    }
}