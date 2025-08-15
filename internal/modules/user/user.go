package user

import "time"

type UserRole string

var (
	RoleUser  UserRole = "USER"
	RoleStaff UserRole = "STAFF"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Phone     *string    `json:"phone"`
	Role      UserRole   `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
