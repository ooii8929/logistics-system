# Environment
- Golang version: go version go1.20.3 darwin/arm64
- Golang Gin Web Framework: gin

# Infra
- Region: 不限
- t2.micro
- Ubuntu 22.04
- Security Group: allow SSH & HTTP

# Features
- 


```
// 建立faker API
// 以 CRUD 的角度確認每張表。e.g.recipient 有沒有要讓他們更新？
//////current_location -> location_id, title, city, address
//////recipient -> id, name, address, phone, created_dt
//// sno, tracking_status, estimated_delivery, recipient, current_location, created_dt, update_dt
//////tracking_history -> record_id, tracking_status, location, created_dt
```




# 本地端開發測試
1. 建立帶有密碼的 redis-server
```
echo "requirepass logisticsredis" > ./redis.conf
redis-server ./redis.conf
redis-cli -h 127.0.0.1 -p 6379
auth logisticsredis
127.0.0.1:6379> auth logisticsredis
OK
```
2. 修改程式碼內的 redis 設定
3. go run .