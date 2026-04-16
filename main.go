package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-crud-api/controller"
	"go-crud-api/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// 1. DB 환경 설정 및 연결
	godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	fmt.Println("✅ DB 연결 완료")

	// ==========================================
	// 2. 의존성 주입 (Dependency Injection)
	// Spring의 @Autowired 과정을 수동으로 연결해 준다.
	// ==========================================
	userRepo := &repository.UserRepository{DB: db}
	userController := &controller.UserController{Repo: userRepo}

	// 3. Gin 라우터(웹 서버) 생성
	router := gin.Default()

	// 헬스 체크
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong! 서버가 정상 작동 중입니다.",
		})
	})

	// ==========================================
	// 4. API 라우팅 (Controller 메서드 매핑)
	// ==========================================
	router.POST("/users", userController.CreateUser)
	router.GET("/users", userController.GetUsers)
	router.GET("/users/:id", userController.GetUser)
	router.PATCH("/users/:id", userController.UpdateUser)
	router.DELETE("/users/:id", userController.DeleteUser)

	// 5. 서버 실행
	fmt.Println("🚀 레이어드 아키텍처 기반 서버가 8080 포트에서 시작됩니다...")
	router.Run(":8080")
}
