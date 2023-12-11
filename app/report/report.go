package report

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
	"time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)


type TrackingSummary struct {
    Status       string `json:"status" gorm:"column:status"`
    TotalRecords int    `json:"total_records" gorm:"column:total_records"`
}

func GetTrackingSummary() (map[string]TrackingSummary, error) {
    username := "root"
    password := "my-secret-password"
    host := "my-mysql"
    port := 3306
    dbName := "logistics"

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }

    query := `
    SELECT status, COUNT(*) as total_records
    FROM tracking_status
    GROUP BY status;
    `

    rows, err := db.Raw(query).Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    summary := make(map[string]TrackingSummary)
    for rows.Next() {
        var s TrackingSummary
        if err := rows.Scan(&s.Status, &s.TotalRecords); err != nil {
            return nil, err
        }
        summary[s.Status] = s
    }

    // 将摘要信息转换为 JSON 数据
    jsonData, err := json.Marshal(summary)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal summary to JSON: %v", err)
    }
    // 获取当前日期和时间
    now := time.Now()

    // 将日期和时间格式化为字符串，例如：2022-01-31T23-30-15
    formattedTime := now.Format("2006-01-02T15-04-05")

    // 使用当前日期和时间创建 JSON 文件名
    jsonFilename := fmt.Sprintf("report-%s.json", formattedTime)

    // 上传 JSON 数据到 S3
    if err := uploadToS3(jsonData, "alvin-report", jsonFilename); err != nil {
        return nil, fmt.Errorf("failed to upload to S3: %v", err)
    }


    return summary, nil
}

func uploadToS3(jsonData []byte, bucket, key string) error {
    // 创建一个新的 AWS session，默认从环境变量、配置文件或 EC2 角色中获取凭据
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-northeast-1")},
    )

    if err != nil {
        return fmt.Errorf("failed to create AWS session: %v", err)
    }

    // 创建一个新的 S3 客户端
    s3Client := s3.New(sess)

    // 设置上传参数
    params := &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        bytes.NewReader(jsonData),
        ContentType: aws.String("application/json"),
    }

    // 执行上传到 S3
    _, err = s3Client.PutObject(params)
    if err != nil {
        return fmt.Errorf("failed to put object to S3: %v", err)
    }

    log.Printf("Successfully uploaded JSON to S3: bucket=%s, key=%s", bucket, key)

    return nil
}