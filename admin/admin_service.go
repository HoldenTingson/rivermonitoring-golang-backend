package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginAdminReq) (*LoginAdminRes, error) {

	a, err := s.Repository.GetAdminByUsername(c, req.Username, req.Password)
	if err != nil {
		return &LoginAdminRes{}, err
	}

	return &LoginAdminRes{
		Id:       strconv.Itoa(int(a.Id)),
		Username: a.Username,
	}, nil
}

func (s *service) GetAdmin(c context.Context, token string) (*Admin, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	a, err := s.Repository.GetAdmin(ctx, token)
	if err != nil {
		return &Admin{}, err
	}

	return &Admin{
		Id:       a.Id,
		Username: a.Username,
		Password: a.Password,
	}, nil
}

func (s *service) AddUser(ctx context.Context, req *CreateUserRequest) error {
	user := User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Language: req.Language,
		Profile:  req.Profile,
	}
	err := s.Repository.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ViewUser(ctx context.Context) (*[]UserResponse, error) {
	var users []UserResponse
	res, err := s.Repository.GetUser(ctx)
	if err != nil {
		return &[]UserResponse{}, err
	}
	for _, u := range *res {
		user := NewUser(u)

		users = append(users, *user)
	}
	return &users, nil
}

func (s *service) ViewUserById(ctx context.Context, id int) (*UserResponse, error) {
	res, err := s.Repository.GetUserById(ctx, id)
	if err != nil {
		return &UserResponse{}, err
	}

	user := UserResponse{
		Id:        res.Id,
		Username:  res.Username,
		Password:  res.Password,
		Email:     res.Email,
		Phone:     res.Phone,
		Language:  res.Language,
		Profile:   res.Profile,
		CreatedAt: res.CreatedAt,
		ChangedAt: res.ChangedAt,
	}

	return &user, nil

}

func NewUser(user User) *UserResponse {
	return &UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Phone:     user.Phone,
		Language:  user.Language,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
		ChangedAt: user.ChangedAt,
	}
}

func (s *service) RemoveUser(ctx context.Context, id int) error {
	err := s.Repository.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangeUser(ctx context.Context, req *UpdateUserRequest, id int) error {
	user := User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Language: req.Language,
		Profile:  req.Profile,
	}

	err := s.Repository.UpdateUser(ctx, &user, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewAllUserCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllUserCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}
