package redis

import (
    "context"
    "fmt"

    "github.com/go-redis/redis/v8"
)

// 返回連線後的redis client, 與錯誤信息
func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password, // 默認為沒有密碼
        DB:       db,       // 使用哪個數據庫，默認為0 
    })

    ctx := context.TODO() // 創建一個context.Context對象，以便與Redis操作結合使用

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("redis connection failed: %v", err)
    }

    fmt.Println("Connected to Redis")
    return rdb, nil
}
