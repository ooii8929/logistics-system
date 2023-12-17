# Offical Golang image as environment
FROM golang:1.20 AS build-env

# 設置工作目錄
WORKDIR /go/src/app

# 將當前目錄下的所有文件複製到容器中的工作目錄
COPY . .

# Constructure golang app
RUN go mod download

# 編譯 Go 應用程序以在 Linux 平台上運行，並禁用 Cgo。將生成的可執行文件命名為 "app"。
RUN CGO_ENABLED=0 GOOS=linux go build -v -o app

# 使用 FROM golang:1.20 AS build-env 來建立 Golang 環境以編譯應用程序，然後使用 FROM scratch 來創建一個僅包含編譯後二進制文件的運行環境
# 最終生成的 Docker 鏡像既有效又具有最小的尺寸
FROM scratch

# 複製構建後的 Go 應用程序到運行的環境中
COPY --from=build-env /go/src/app/app /app

# 運行 Go 應用程序
CMD ["/app"]
