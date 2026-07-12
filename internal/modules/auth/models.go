package auth

import "time"

type User struct {
	UserID      string    `json:"user_id" db:"user_id"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Admin struct {
	AdminID     string    `json:"admin_id" db:"admin_id"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
