# 1. Hướng Dẫn Chạy Ứng Dụng
## Các Bước Cài Đặt

- Bước 1: Cài Đặt Go.   
Truy cập trang chính thức của Go để tải và cài đặt: https://golang.org/dl/

    Sau khi cài đặt xong, kiểm tra lại bằng lệnh sau trong terminal:
    ```bash
    go version
    ```
    Nếu cài đặt thành công, bạn sẽ thấy phiên bản Go hiện tại.

- Bước 2: Cấu hình dự án.   
    Chạy các câu lệnh dưới đây để khởi chạy project
    ```bash
    go mod tidy
    go build -o main; ./main
    ```
- Bước 3: Project sẽ listen port 8080 và chạy mode TLS.   
    ```[GIN-debug] Listening and serving HTTPS on :8080```