package user

import (
	"context"
)

type User struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Language  string `json:"language" db:"language"`
	Profile   string `json:"profile" db:"profile"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUser(ctx context.Context, token string) (*User, error)
	UpdatePassword(ctx context.Context, email string, newPassword string) error
	UpdateProfile(ctx context.Context, user *User, id int, update string) error
}

type CreateUserReq struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Language  string `json:"language" db:"language"`
	Profile   string `json:"profile" db:"profile"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type CreateUserRes struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type ChangePasswordReq struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type LoginUserRes struct {
	accessToken string
	Id          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	GetUser(c context.Context, token string) (*User, error)
	ChangePassword(c context.Context, req *ChangePasswordReq) error
	ChangeProfile(ctx context.Context, req *UpdateUserRequest, id int, update string) error
}
