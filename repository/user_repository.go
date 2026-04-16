package repository

import (
	"database/sql"
	"go-crud-api/models"
)

// UserRepository 구조체 (Spring의 @Repository 클래스 역할)
type UserRepository struct {
	DB *sql.DB // DB 연결 객체를 주입받아 사용한다
}

// 1. 유저 생성 (Create)
func (r *UserRepository) CreateUser(user *models.User) error {
	result, err := r.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

// 2. 전체 유저 조회 (Read All)
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// 3. 단일 유저 조회 (Read One)
func (r *UserRepository) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	return user, err
}

// 4. 유저 정보 수정 (Update - PATCH)
func (r *UserRepository) UpdateUser(user models.User) error {
	_, err := r.DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, user.ID)
	return err
}

// 5. 유저 삭제 (Delete)
func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
