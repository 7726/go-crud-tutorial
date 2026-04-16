package controller

import (
	"go-crud-api/models"
	"go-crud-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 구조체 (Spring의 @RestController 역할)
type UserController struct {
	Repo *repository.UserRepository // Repository 의존성 주입 (DI)
}

// 1. 유저 생성 (Create)
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	// Repository 계층 호출
	if err := c.Repo.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB 저장 실패"})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// 2. 전체 유저 조회 (Read All)
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.Repo.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "DB 조회 실패"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// 3. 단일 유저 조회 (Read One)
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.Repo.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "유저를 찾을 수 없습니다."})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// 4. 유저 정보 수정 (Update - PATCH)
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	type UpdateReq struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}
	var req UpdateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다."})
		return
	}

	// 기존 유저 조회
	user, err := c.Repo.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "수정할 유저를 찾을 수 없습니다."})
		return
	}

	// Dirty Checking 처럼 값이 들어온 부분만 병합
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := c.Repo.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "수정 실패"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "부분 수정 완료", "updated_user": user})
}

// 5. 유저 삭제 (Delete)
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.Repo.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "삭제 실패"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "삭제 완료"})
}
