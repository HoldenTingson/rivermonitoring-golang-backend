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
	Profile   string `json:"profile" db:"profile"`
	CreatedAt string `json:"created_at" db:"created_at"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}
type CreateUserRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Profile  string `json:"profile" db:"profile"`
}

type UserResponse struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Profile   string `json:"profile" db:"profile"`
	CreatedAt string `json:"created_at" db:"created_at"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type Repository interface {
	RegisterUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUser(ctx context.Context, token string) (*User, error)
	CheckPassword(ctx context.Context, email string, currentPassword string) (bool, error)
	UpdatePassword(ctx context.Context, email string, newPassword string) error
	UpdateProfile(ctx context.Context, user *User, id int, update string) error
	SendEmail(ctx context.Context, email string, hash string) (string, error)
	CreatePassword(ctx context.Context, username string, hash string, password string) error
	GetUsers(ctx context.Context) (*[]User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) error
}

type Service interface {
	RegisterUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	GetUser(c context.Context, token string) (*User, error)
	ChangePassword(c context.Context, req *ChangePasswordReq) error
	ChangeProfile(ctx context.Context, req *UpdateUserRequest, id int, update string) error
	SendEmail(ctx context.Context, email string) error
	CreatePassword(ctx context.Context, req *CreatePasswordReq) error
	ViewUser(ctx context.Context) (*[]UserResponse, error)
	ViewUserById(ctx context.Context, id int) (*UserResponse, error)
	AddUser(ctx context.Context, req *CreateUserRequest) error
	RemoveUser(ctx context.Context, id int) error
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

type CreatePasswordReq struct {
	Username         string `json:"username"`
	EmailVerPassword string `json:"emailVerPassword"`
	Password         string `json:"password"`
}

type LoginUserRes struct {
	accessToken string
	Id          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
}

type SendEmailReq struct {
	Email string `json:"email"`
}
