package report

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "crypto/tls"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "gorm.io/gorm"
)

type TrackingSummary struct {
    TrackingStatus string `json:"tracking_status" gorm:"column:tracking_status"`
    TotalRecords   int    `json:"total_records" gorm:"column:total_records"`
}

func GetTrackingSummary(db *gorm.DB) (map[string]TrackingSummary, error) {
    log.Println("Starting GetTrackingSummary")

    query := `
    SELECT tracking_status, COUNT(*) as total_records
    FROM trackings
    GROUP BY tracking_status;
    `

    rows, err := db.Raw(query).Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    summary := make(map[string]TrackingSummary)
    for rows.Next() {
        var s TrackingSummary
        if err := rows.Scan(&s.TrackingStatus, &s.TotalRecords); err != nil {
            return nil, err
        }
        summary[s.TrackingStatus] = s
    }

    log.Println("Summary data:", summary)

    jsonData, err := json.Marshal(summary)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal summary to JSON: %v", err)
    }
    now := time.Now()

    formattedTime := now.Format("2006-01-02T15-04-05")

    jsonFilename := fmt.Sprintf("report-%s.json", formattedTime)

    log.Printf("Uploading %s to S3 bucket\n", jsonFilename)

    if err := uploadToS3(jsonData, "alvin-report", jsonFilename); err != nil {
        return nil, fmt.Errorf("failed to upload to S3: %v", err)
    }

    log.Println("GetTrackingSummary successfully completed")
    return summary, nil
}

func uploadToS3(jsonData []byte, bucket, key string) error {
    log.Println("Starting uploadToS3")

    sess, err := session.NewSessionWithOptions(session.Options{
        Config: aws.Config{
            Region: aws.String("ap-northeast-1"),
            S3Disable100Continue: aws.Bool(true),
            HTTPClient: &http.Client{
                Transport: &http.Transport{
                    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                },
            },
            CredentialsChainVerboseErrors: aws.Bool(true),
        },
        // // 使用 SharedConfigEnable 和 Profile 配置
        // SharedConfigState: session.SharedConfigEnable,
        // Profile:           "alvin",
    })

    if err != nil {
        return fmt.Errorf("failed to create AWS session: %v", err)
    }

    s3Client := s3.New(sess)

    params := &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        bytes.NewReader(jsonData),
        ContentType: aws.String("application/json"),
    }

    log.Println("Uploading to S3 bucket:", bucket)
    _, err = s3Client.PutObject(params)
    if err != nil {
        return fmt.Errorf("failed to put object to S3: %v", err)
    }

    log.Println("Upload to S3 completed successfully")
    return nil
}
