package user

import (
	"BanjirEWS/util"
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

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		Id:       strconv.Itoa(int(r.Id)),
		Username: r.Username,
		Email:    r.Email,
	}
	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.Id)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{
		accessToken: ss,
		Id:          strconv.Itoa(int(u.Id)),
		Username:    u.Username,
	}, nil
}

func (s *service) GetUser(c context.Context, token string) (*User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUser(ctx, token)
	if err != nil {
		return &User{}, err
	}

	return &User{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
		Language: u.Language,
		Profile:  u.Profile,
	}, nil
}

func (s *service) ChangeProfile(ctx context.Context, req *UpdateUserRequest, id int, update string) error {
	user := User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Phone:     req.Phone,
		Language:  req.Language,
		Profile:   req.Profile,
		ChangedAt: req.ChangedAt,
	}

	err := s.Repository.UpdateProfile(ctx, &user, id, update)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangePassword(c context.Context, req *ChangePasswordReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Generate a hashed password for the new password
	newHashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update the user's password in the repository
	err = s.Repository.UpdatePassword(ctx, req.Email, newHashedPassword)
	if err != nil {
		return err
	}

	return nil
}
