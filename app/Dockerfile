# 使用官方 Golang 鏡像作為構建環境
FROM golang:1.20 AS build-env

# 設置工作目錄
WORKDIR /go/src/app

# 將你的go應用程序源代碼複制到容器中
COPY . .

# 構建 Go 應用程序
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

# 使用刮條（Scratch）鏡像作為運行環境
FROM scratch

# 複製構建後的 Go 應用程序到這個環境中
COPY --from=build-env /go/src/app/app /app

# 運行 Go 應用程序
CMD ["/app"]
