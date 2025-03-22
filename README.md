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

# 2. Tài liệu API
### ***Do thời gian có hạn nên tôi chưa thể hoàn thành hết toàn bộ yêu cầu***
```
Xác thực tĩnh thông qua token tại header
Authorization = 1234567890abcdefjustforspeed
```
## 1. API lấy dữ liệu
route: /api/custom_search/  
method: POST    
body:
```
{
    "filters": {
    },
    "get_fields": [],
    "table": "products",
    "limit": 5,
    "offset": 1
}
```
giải thích:
 - filters: điều kiện tìm kiếm
 - get_fields: danh sách trường thông tin cần truy xuất, **nếu không truyền gì mặc định lấy tất cả các trường**
 - table: chỉ định bảng dữ liệu tìm kiếm
 - limit: số lượng bản ghi trả về
 - offset: vị trí dữ liệu

body example:
```
{
    "filters": {
        "reference": "PROD-202401-003",
        "status": "Available"
    },
    "get_fields": [],
    "table": "products",
    "limit": 5,
    "offset": 1
}
```

Dữ liệu trả về:
```
{
    "data": [
        {
            "id": "ab15c342-6903-4f2f-abb8-de65d3cdd847",
            "reference": "PROD-202401-003",
            "name": "Sofa Deluxe",
            "added_date": "2024-01-03T00:00:00Z",
            "status": "Available",
            "category_id": "e8aa616d-dacb-4924-a270-b77030af1a02",
            "price": 899.99,
            "stock_city": "Marseille",
            "supplier_id": "7ee42cdd-26db-4649-8cea-49c50d23fefb",
            "quantity": 20
        }
    ],
    "message": "success",
    "metadata": {
        "limit": 5,
        "page": 0
    },
    "status": true
}
```

## 2. API lấy một bản ghi trong bảng suppliers
route: /api/suppliers/<:id>  
method: GET     
Dữ liệu trả về:
```
{
    "data": [
        {
            "id": "7ee42cdd-26db-4649-8cea-49c50d23fefb",
            "name": "Home Goods"
        }
    ],
    "message": "success",
    "metadata": null,
    "status": true
}
```

## 3. API lấy một bản ghi trong bảng products
route: /api/products/<:id>  
method: GET     
Dữ liệu trả về:
```
{
    "data": [
        {
            "id": "debd60f6-eeb4-4e52-b5a8-26cda75caec1",
            "reference": "PROD-202401-001",
            "name": "Smartphone YX",
            "added_date": "2024-01-01T00:00:00Z",
            "status": "Available",
            "category_id": "94d0da61-0bbe-4be8-8435-2b72f03a29ea",
            "price": 10,
            "stock_city": "Paris",
            "supplier_id": "4f8ce93f-46c2-4d20-8a27-92fdbf6ee464",
            "quantity": 50
        }
    ],
    "message": "success",
    "metadata": null,
    "status": true
}
```

## 4. API lấy một bản ghi trong bảng categories
route: /api/categories/<:id>  
method: GET     
Dữ liệu trả về:
```
{
    "data": [
        {
            "id": "62b18bbf-d94d-4768-a072-23992928b59b",
            "name": "Books"
        }
    ],
    "message": "success",
    "metadata": null,
    "status": true
}
```

## 5. API cập nhật bản ghi
route: /api/<:tablename>/<:id>  
method: POST    
Body example:
```
{
    "price": 10
}
```
Dữ liệu trả về:
```
{
    "data": null,
    "message": "success",
    "metadata": null,
    "status": true
}
```