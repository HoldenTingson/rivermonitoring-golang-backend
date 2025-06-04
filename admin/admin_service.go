package admin

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

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) CreateAdmin(c context.Context, req *CreateAdminReq) (*CreateAdminRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	a := &Admin{
		Username: req.Username,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateAdmin(ctx, a)
	if err != nil {
		return nil, err
	}

	res := &CreateAdminRes{
		Id:       r.Id,
		Username: r.Username,
	}
	return res, nil
}

func (s *service) Login(c context.Context, req *LoginAdminReq) (*LoginAdminRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	a, err := s.Repository.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		return &LoginAdminRes{}, err
	}

	err = util.CheckPassword(req.Password, a.Password)
	if err != nil {
		return &LoginAdminRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(a.Id)),
		Username: a.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(a.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginAdminRes{}, err
	}

	return &LoginAdminRes{
		accessToken: ss,
		Id:          strconv.Itoa(int(a.Id)),
		Username:    a.Username,
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

func (s *service) ViewAdmin(ctx context.Context) (*[]AdminResponse, error) {
	var admins []AdminResponse
	res, err := s.Repository.GetAdmins(ctx)
	if err != nil {
		return &[]AdminResponse{}, err
	}
	for _, a := range *res {
		admin := NewAdmin(a)

		admins = append(admins, *admin)
	}
	return &admins, nil
}

func (s *service) ViewAdminById(ctx context.Context, id int) (*AdminResponse, error) {
	res, err := s.Repository.GetAdminById(ctx, id)
	if err != nil {
		return &AdminResponse{}, err
	}

	admin := AdminResponse{
		Id:        int64(res.Id),
		Username:  res.Username,
		CreatedAt: res.CreatedAt,
	}

	return &admin, nil

}

func NewAdmin(admin Admin) *AdminResponse {
	return &AdminResponse{
		Id:        int64(admin.Id),
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt,
	}
}

func (s *service) RemoveAdmin(ctx context.Context, id int) error {
	err := s.Repository.DeleteAdmin(ctx, id)
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
