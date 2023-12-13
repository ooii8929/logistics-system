package database

import (
    "log"
    "fmt"
    "time"

    // Fake API
    "math/rand"
    
    "gorm.io/gorm"
    _ "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "logistics-track/models"
)

// 進行連接，回傳 DB 連線
func ConnectToMySQL(username, password, host string, port int) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        return nil, err
    }

    return db, nil
}

// 進行連接，回傳 DB 連線
func ConnectToMySQLwithTable(username, password, host, dbName string, port int) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        return nil, err
    }

    return db, nil
}


// 建立 DB，並初始化基礎數據
func CreateAndInitializeDatabase(username, password, host string, port int, databaseName string) (error) {

    // Create a new database
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("Error connecting to MySQL while initializing the database: %v", err)
    }

    err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", databaseName)).Error
    if err != nil {
        return fmt.Errorf("Error creating the '%s' database: %v", databaseName, err)
    }
    fmt.Printf("Database '%s' created successfully.\n", databaseName)

    // Reconnect with the new database
    dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, databaseName)
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("Error reconnecting to MySQL with the new '%s' database: %v", databaseName, err)
    }

    // Create table
    CreateTables(db)
    db.AutoMigrate(&models.Location{}, &models.Recipient{}, &models.Tracking{}, &models.TrackingHistory{})

    // Insert data
    initLocations(db)
    initRecipients(db)

    return nil
}


func CreateTables(db *gorm.DB) {
    tables := []string{
        `
        CREATE TABLE IF NOT EXISTS locations (
            location_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            title TEXT NOT NULL,
            city TEXT NOT NULL,
            address TEXT NOT NULL,
            created_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )`,
        `
        CREATE TABLE IF NOT EXISTS recipients (
            id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
            name TEXT NOT NULL,
            address TEXT NOT NULL,
            phone TEXT NOT NULL,
            created_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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


func initLocations(db *gorm.DB) {
    locations := []map[string]interface{}{
        {"location_id": 7, "title": "台北物流中⼼", "city": "台北市", "address": "台北市中正區忠孝東路100號"},
        {"location_id": 13, "title": "新⽵物流中⼼", "city": "新⽵市", "address": "新⽵市東區光復路⼀段101號"},
        {"location_id": 24, "title": "台中物流中⼼", "city": "台中市", "address": "台中市⻄區⺠⽣路200號"},
        {"location_id": 3, "title": "桃園物流中⼼", "city": "桃園市", "address": "桃園市中壢區中央⻄路三段150號"},
        {"location_id": 18, "title": "⾼雄物流中⼼", "city": "⾼雄市", "address": "⾼雄市前⾦區成功⼀路82號"},
        {"location_id": 9, "title": "彰化物流中⼼", "city": "彰化市", "address": "彰化市中⼭路⼆段250號"},
        {"location_id": 15, "title": "嘉義物流中⼼", "city": "嘉義市", "address": "嘉義市東區⺠族路380號"},
        {"location_id": 6, "title": "宜蘭物流中⼼", "city": "宜蘭市", "address": "宜蘭市中⼭路⼆段58號"},
        {"location_id": 21, "title": "屏東物流中⼼", "city": "屏東市", "address": "屏東市⺠⽣路300號"},
        {"location_id": 1, "title": "花蓮物流中⼼", "city": "花蓮市", "address": "花蓮市國聯⼀路100號"},
        {"location_id": 4, "title": "台南物流中⼼", "city": "台南市", "address": "台南市安平區建平路18號"},
        {"location_id": 11, "title": "南投物流中⼼", "city": "南投市", "address": "南投市⾃由路67號"},
        {"location_id": 23, "title": "雲林物流中⼼", "city": "雲林市", "address": "雲林市中正路五段120號"},
        {"location_id": 14, "title": "基隆物流中⼼", "city": "基隆市", "address": "基隆市信⼀路50號"},
        {"location_id": 8, "title": "澎湖物流中⼼", "city": "澎湖縣", "address": "澎湖縣⾺公市中正路200號"},
        {"location_id": 19, "title": "⾦⾨物流中⼼", "city": "⾦⾨縣", "address": "⾦⾨縣⾦城鎮⺠⽣路90號"},
    }

    now := time.Now()

    for _, location := range locations {
        res := db.Exec("INSERT IGNORE INTO locations (location_id, title, city, address, created_dt) VALUES (?, ?, ?, ?, ?)",
            location["location_id"], location["title"], location["city"], location["address"], now)
        if err := res.Error; err != nil {
            log.Fatalf("Failed to insert location: %v", err)
        }
    }
}

func initRecipients(db *gorm.DB) {
    recipients := []map[string]interface{}{
        {"id": 1234, "name": "賴⼩賴", "address": "台北市中正區仁愛路⼆段99號", "phone": "091234567"},
        {"id": 1235, "name": "陳⼤明", "address": "新北市板橋區⽂化路⼀段100號", "phone": "092345678"},
        {"id": 1236, "name": "林⼩芳", "address": "台中市⻄區⺠⽣路200號", "phone": "093456789"},
        {"id": 1237, "name": "張美玲", "address": "⾼雄市前⾦區成功⼀路82號", "phone": "094567890"},
        {"id": 1238, "name": "王⼩明", "address": "台南市安平區建平路18號", "phone": "095678901"},
        {"id": 1239, "name": "劉⼤華", "address": "新⽵市東區光復路⼀段101號", "phone": "096789012"},
        {"id": 1240, "name": "⿈⼩琳", "address": "彰化市中⼭路⼆段250號", "phone": "097890123"},
        {"id": 1241, "name": "吳美美", "address": "花蓮市國聯⼀路100號", "phone": "098901234"},
        {"id": 1242, "name": "蔡⼩虎", "address": "屏東市⺠⽣路300號", "phone": "099012345"},
        {"id": 1243, "name": "鄭⼤勇", "address": "基隆市信⼀路50號", "phone": "091123456"},
        {"id": 1244, "name": "謝⼩珍", "address": "嘉義市東區⺠族路380號", "phone": "092234567"},
        {"id": 1245, "name": "潘⼤為", "address": "宜蘭市中⼭路⼆段58號", "phone": "093345678"},
        {"id": 1246, "name": "趙⼩梅", "address": "南投市⾃由路67號", "phone": "094456789"},
        {"id": 1247, "name": "周⼩⿓", "address": "雲林市中正路五段120號", "phone": "095567890"},
        {"id": 1248, "name": "李⼤同", "address": "澎湖縣⾺公市中正路200號", "phone": "096678901"},
        {"id": 1249, "name": "陳⼩凡", "address": "⾦⾨縣⾦城鎮⺠⽣路90號", "phone": "097789012"},
        {"id": 1250, "name": "楊⼤明", "address": "台北市信義區松仁路50號", "phone": "098890123"},
        {"id": 1251, "name": "吳⼩雯", "address": "新北市中和區景平路100號", "phone": "099901234"},
    }

    now := time.Now()

    for _, recipient := range recipients {
        res := db.Exec("INSERT IGNORE INTO recipients (id, name, address, phone, created_dt) VALUES (?, ?, ?, ?, ?)",
            recipient["id"], recipient["name"], recipient["address"], recipient["phone"], now)
        if err := res.Error; err != nil {
            log.Fatalf("Failed to insert recipient: %v", err)
        }
    }

}



func GenerateFakeData(db *gorm.DB, num int) []models.Tracking {

    allStatuses := []models.TrackingStatus{
        models.Created,
        models.PackageReceived,
        models.InTransit,
        models.OutForDelivery,
        models.DeliveryAttempted,
        models.Delivered,
        models.ReturnedToSender,
        models.Exception,
    }    
    trackingList := make([]models.Tracking, 0)

    var locationIDs []int
    db.Table("locations").Select("location_id").Find(&locationIDs)

    var recipientIDs []int
    db.Table("recipients").Select("id").Find(&recipientIDs)

    if num > 0 {

        // Insert fake trackings
        for i := 0; i < num; i++ {
            tr := models.Tracking{
                TrackingStatus:    string(allStatuses[rand.Intn(len(allStatuses))]),
                RecipientId:       recipientIDs[rand.Intn(len(recipientIDs))],
                CurrentLocation:   locationIDs[rand.Intn(len(locationIDs))],
                EstimatedDelivery: time.Now().Add(time.Duration(rand.Intn(168)) * time.Hour),
                CreatedDt:         time.Now(),
                UpdatedDt:         time.Now(),
            }
            db.Create(&tr)
            trackingList = append(trackingList, tr)
        }

        // Insert fake tracking histories
        for i := 1; i < num; i++ {
            randomTrackingIndex := rand.Intn(len(trackingList))
            trhist := models.TrackingHistory{
                TrackingStatus: string(allStatuses[rand.Intn(len(allStatuses))]),
                Sno:            trackingList[randomTrackingIndex].Sno,
                LocationID:     locationIDs[rand.Intn(len(locationIDs))],
                CreatedDt:      time.Now(),
            }
            db.Create(&trhist)
        }
    }
    return trackingList
}
