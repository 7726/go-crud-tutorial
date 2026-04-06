package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// 1. .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env 파일을 찾을 수 없습니다.")
	}

	// 2. 환경 변수 읽기
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 3. DB 연결 문자열(DSN) 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// 4. DB 연결 시도
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 5. 실제 연결 확인 (Ping)
	err = db.Ping()
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	fmt.Println("✅ MySQL 데이터베이스 연결 성공!")
}
