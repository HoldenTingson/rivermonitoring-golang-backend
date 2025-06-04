package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var u CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.RegisterUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var user LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    u.accessToken,
		Path:     "/",
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	res := &LoginUserRes{
		Username: u.Username,
		Id:       u.Id,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUser(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid JWT token"})
		return
	}
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid JWT token"})
		return
	}
	claims := token.Claims.(*jwt.RegisteredClaims)

	u, err := h.Service.GetUser(c.Request.Context(), claims.Issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &User{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Phone:    u.Phone,
		Profile:  u.Profile,
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) ChangeProfile(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	update := c.Query("update")

	err := h.Service.ChangeProfile(c.Request.Context(), &req, id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Change User")
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout Succesful"})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.ChangePassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *Handler) SendEmail(c *gin.Context) {
	var req SendEmailReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.SendEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})

}

func (h *Handler) CreatePassword(c *gin.Context) {
	var req CreatePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.CreatePassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password created successfully"})
}

func (h *Handler) AddUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.AddUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully Add User")
}

func (h *Handler) ViewUser(c *gin.Context) {
	res, err := h.Service.ViewUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
func (h *Handler) ViewUserById(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)
	res, err := h.Service.ViewUserById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) RemoveUser(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.RemoveUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Remove User")
}
