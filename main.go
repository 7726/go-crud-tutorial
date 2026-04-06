package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// User 구조체 (Spring의 Entity / DTO 역할)
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB // 전역 변수로 DB 핸들러 유지

func main() {
	// 1. DB 환경 설정 및 연결 (Step 4와 동일)
	godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	fmt.Println("✅ DB 연결 완료")

	// 2. Gin 라우터(웹 서버) 생성
	router := gin.Default()

	// 3. 테스트용 API (Spring의 @GetMapping("/ping"))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong! 서버가 정상 작동 중입니다.",
		})
	})

	// 4. 서버 실행 (기본 8080 포트)
	fmt.Println("🚀 서버가 8080 포트에서 시작됩니다...")

	// ==========================================
	// 여기서부터 CRUD API (Spring의 Controller 역할)
	// ==========================================

	// 1. Create (유저 생성 - POST)
	router.POST("/users", func(c *gin.Context) {
		var user User
		// @RequestBody처럼 JSON 데이터를 User 구조체에 바인딩
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
			return
		}

		// PreparedStatement 처럼 '?' 를 사용하여 SQL Injection 방지
		result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 저장 실패"})
			return
		}

		id, _ := result.LastInsertId()
		user.ID = int(id)
		c.JSON(http.StatusCreated, user) // 201 Created 응답
	})

	// 2. Read All (전체 유저 조회 - GET)
	router.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, email FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 조회 실패"})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				continue
			}
			users = append(users, u)
		}
		c.JSON(http.StatusOK, users)
	})

	// 3. Read One (단일 유저 조회 - GET)
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id") // @PathVariable 역할
		var user User

		err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "유저를 찾을 수 없습니다."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 에러"})
			}
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// 4. Update (유저 부분 수정 - PATCH)
	router.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 1) 부분 수정을 위해 값이 들어왔는지 확인할 임시 구조체 (포인터 사용)
		type UpdateReq struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}
		var req UpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
			return
		}

		// 2) 기존 데이터 조회 (먼저 DB에서 타겟 유저를 가져옵니다)
		var user User
		err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "수정할 유저를 찾을 수 없습니다."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 조회 에러"})
			}
			return
		}

		// 3) 값이 들어온 필드만 덮어쓰기 (Spring의 Dirty Checking 또는 부분 병합과 유사)
		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Email != nil {
			user.Email = *req.Email
		}

		// 4) DB에 최종 업데이트 반영
		_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "수정 실패"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "부분 수정 완료",
			"updated_user": user,
		})
	})

	// 5. Delete (유저 삭제 - DELETE)
	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "삭제 실패"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "삭제 완료"})
	})

	// ==========================================
	// CRUD API 끝
	// ==========================================

	router.Run(":8080")
}
