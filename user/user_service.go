package user

import (
	"BanjirEWS/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

func (s *service) RegisterUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
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

	r, err := s.Repository.RegisterUser(ctx, u)
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
		Profile:  u.Profile,
	}, nil
}

func (s *service) ChangeProfile(ctx context.Context, req *UpdateUserRequest, id int, update string) error {
	user := User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Phone:     req.Phone,
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

	same, err := s.Repository.CheckPassword(ctx, req.Email, req.CurrentPassword)
	if err != nil {
		return err
	}

	if !same {
		return errors.New("invalid current password")
	}

	newHashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.Repository.UpdatePassword(ctx, req.Email, newHashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendEmail(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}

	emailVerPassword := string(emailVerRandRune)
	var emailVerPasswordHash []byte
	emailVerPasswordHash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	emailVerPasswordHashStr := string(emailVerPasswordHash)

	username, err := s.Repository.SendEmail(ctx, email, emailVerPasswordHashStr)
	if err != nil {
		return err
	}

	from := "holdentingson33@gmail.com"
	password := "mxhuhizmirbrqzli"
	to := []string{email}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := "Subject: Go Banjir Account Recovery\n"
	body := "<body><a rel=\"nofollow noopener noreferrer\" target=\"_blank\" href=\"https://gobanjirclient.netlify.app/createPassword?u=" + username + "&evpw=" + emailVerPassword + "\">Change Password</a></body>"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + body)
	auth := smtp.PlainAuth("", from, password, host)
	err = smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreatePassword(ctx context.Context, req *CreatePasswordReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	err = s.Repository.CreatePassword(ctx, req.Username, req.EmailVerPassword, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func NewUser(user User) *UserResponse {
	return &UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Phone:     user.Phone,
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
		Profile:   res.Profile,
		CreatedAt: res.CreatedAt,
		ChangedAt: res.ChangedAt,
	}

	return &user, nil

}

func (s *service) ViewUser(ctx context.Context) (*[]UserResponse, error) {
	var users []UserResponse
	res, err := s.Repository.GetUsers(ctx)
	if err != nil {
		return &[]UserResponse{}, err
	}
	for _, u := range *res {
		user := NewUser(u)

		users = append(users, *user)
	}
	return &users, nil
}

func (s *service) AddUser(ctx context.Context, req *CreateUserRequest) error {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Profile:  req.Profile,
	}

	err = s.Repository.CreateUser(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
