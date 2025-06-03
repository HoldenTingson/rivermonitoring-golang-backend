package admin

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

func (h *Handler) CreateAdmin(c *gin.Context) {
	var a CreateAdminReq
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateAdmin(c.Request.Context(), &a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var admin LoginAdminReq
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	a, err := h.Service.Login(c.Request.Context(), &admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", a.accessToken, 3600, "/", "localhost", false, true)
	res := &LoginAdminRes{
		Username: a.Username,
		Id:       a.Id,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetAdmin(c *gin.Context) {
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

	a, err := h.Service.GetAdmin(c.Request.Context(), claims.Issuer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &Admin{
		Id:       a.Id,
		Username: a.Username,
		Password: a.Password,
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout Succesful"})
}

func (h *Handler) UploadImage(c *gin.Context) {
	var req UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	base64Data := req.Blob[strings.IndexByte(req.Blob, ',')+1:]

	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode base64 data"})
		return
	}

	filePath := filepath.Join(req.Path, req.Filename)

	err = ioutil.WriteFile(filePath, decodedData, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func (h *Handler) ViewAdmin(c *gin.Context) {
	res, err := h.Service.ViewAdmin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ViewAdminById(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)
	res, err := h.Service.ViewAdminById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) RemoveAdmin(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.RemoveAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Remove Admin")
}

func (h *Handler) DisplayUserCount(c *gin.Context) {
	res, err := h.Service.ViewAllUserCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
