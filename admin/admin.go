package admin

import (
	"context"
)

type Admin struct {
	Id        int    `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type CreateAdminReq struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type CreateAdminRes struct {
	Id       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
}

type LoginAdminReq struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type LoginAdminRes struct {
	accessToken string
	Id          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
}

type UpdateUserRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Profile  string `json:"profile" db:"profile"`
}

type AdminResponse struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	CreatedAt string `json:"created_at" db:"created_at"`
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
	Profile   string `json:"profile" db:"profile"`
	CreatedAt string `json:"created_at" db:"created_at"`
	ChangedAt string `json:"changed_at" db:"changed_at"`
}

type Repository interface {
	CreateAdmin(ctx context.Context, admin *Admin) (*Admin, error)
	GetAdmin(ctx context.Context, token string) (*Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*Admin, error)
	GetAdminById(ctx context.Context, id int) (*Admin, error)
	DeleteAdmin(ctx context.Context, id int) error
	GetAdmins(ctx context.Context) (*[]Admin, error)
	GetAllUserCount(ctx context.Context) (int, error)
}

type Service interface {
	CreateAdmin(ctx context.Context, req *CreateAdminReq) (*CreateAdminRes, error)
	GetAdmin(c context.Context, token string) (*Admin, error)
	Login(c context.Context, req *LoginAdminReq) (*LoginAdminRes, error)
	RemoveAdmin(ctx context.Context, id int) error
	ViewAdmin(ctx context.Context) (*[]AdminResponse, error)
	ViewAdminById(ctx context.Context, id int) (*AdminResponse, error)
	ViewAllUserCount(ctx context.Context) (int, error)
}
