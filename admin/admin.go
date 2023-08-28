package admin

import (
	"context"
)

type Admin struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type LoginAdminReq struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type LoginAdminRes struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}

type UpdateUserRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Language string `json:"language" db:"language"`
	Profile  string `json:"profile" db:"profile"`
}

type CreateUserRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Language string `json:"language" db:"language"`
	Profile  string `json:"profile" db:"profile"`
}

type UserResponse struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Language  string `json:"language" db:"language"`
	Profile   string `json:"profile" db:"profile"`
	CreatedAt string `json:"created_at" db:"created_at"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type UploadRequest struct {
	Blob     string `json:"blob"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type User struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Language  string `json:"language" db:"language"`
	Profile   string `json:"profile" db:"profile"`
	CreatedAt string `json:"created_at" db:"created_at"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type Repository interface {
	GetAdmin(ctx context.Context, token string) (*Admin, error)
	GetAdminByUsername(ctx context.Context, username string, password string) (*Admin, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User, id int) error
	DeleteUser(ctx context.Context, id int) error
	GetUser(ctx context.Context) (*[]User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
	GetAllUserCount(ctx context.Context) (int, error)
}

type Service interface {
	GetAdmin(c context.Context, token string) (*Admin, error)
	Login(c context.Context, req *LoginAdminReq) (*LoginAdminRes, error)
	AddUser(ctx context.Context, req *CreateUserRequest) error
	ChangeUser(ctx context.Context, req *UpdateUserRequest, id int) error
	RemoveUser(ctx context.Context, id int) error
	ViewUser(ctx context.Context) (*[]UserResponse, error)
	ViewUserById(ctx context.Context, id int) (*UserResponse, error)
	ViewAllUserCount(ctx context.Context) (int, error)
}
