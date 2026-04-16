# 🚀 Go-Sprint CRUD API

단기간에 Golang 생태계를 파악하고 실무 환경에 적응하기 위해 구축한 초경량 CRUD API 서버입니다.

## 🛠 Tech Stack
- **Language**: Go (v1.25+)
- **Web Framework**: Gin (`github.com/gin-gonic/gin`)
- **Database**: MySQL 8.0+
- **Driver/Config**: `go-sql-driver/mysql`, `joho/godotenv`

## 💡 Key Features & Architecture Notes
- **RESTful API**: `GET`, `POST`, `PATCH`(부분 수정), `DELETE` 메서드를 활용한 표준 규격 준수.
- **Security & Safety**: `?` 파라미터 바인딩을 통한 SQL Injection 방지 및 `*`(포인터)를 활용한 Null/Empty-string 안전성 확보.
- **Configuration**: `.env` 파일을 통한 민감 정보 분리.

## 🚀 Quick Start

### 1. Database Setup
```sql
CREATE DATABASE go_study;
USE go_study;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL
);
```

### 2. Environment Variables
프로젝트 루트에 `.env` 파일을 생성하고 아래 양식에 맞게 작성합니다.
```env
DB_USER=your_id
DB_PASS=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=go_study
```

### 3. Run Server
```bash
go mod tidy
go run main.go
```
서버는 기본적으로 `http://localhost:8080` 에서 실행됩니다.

## 📡 API Endpoints

| Method | Endpoint | Description | Request Body (Example) |
|--------|----------|-------------|-------------------------|
| GET | `/ping` | Health Check | - |
| POST | `/users` | 유저 생성 | `{"name":"홍길동", "email":"abc@test.com"}` |
| GET | `/users` | 전체 유저 조회 | - |
| GET | `/users/:id` | 단일 유저 조회 | - |
| PATCH | `/users/:id` | 유저 부분 수정 | `{"name":"홍길동_수정됨"}` (선택적 필드) |
| DELETE| `/users/:id` | 유저 삭제 | - |