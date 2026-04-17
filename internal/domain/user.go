package domain

import (
	"time"

	"github.com/lib/pq"
)

type Role string

const (
	RoleAdmin   Role = "ADMIN"
	RoleTeacher Role = "TEACHER"
	RoleStudent Role = "STUDENT"
)

type User struct {
	ID        string         `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Password  string         `json:"-" db:"password"`
	Roles     pq.StringArray `json:"roles" db:"roles"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
}

type UserProfile struct {
	Name   string   `json:"name" example:"Иван Иванов"`
	Email  string   `json:"email" example:"ivan@example.com"`
	Roles  []string `json:"roles" example:"STUDENT"`
	Avatar *string  `json:"avatar" example:"https://storage.yandexcloud.net/bucket/avatars/user.webp"`
}
