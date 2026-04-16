package models

// User 구조체 (Spring의 Entity / DTO 역할)
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
