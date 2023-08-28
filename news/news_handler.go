package news

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) AddNews(c *gin.Context) {
	var r CreateNewsRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.AddNews(c.Request.Context(), &r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) RemoveNews(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.RemoveNews(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Remove News")
}

func (h *Handler) ChangeNews(c *gin.Context) {
	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.ChangeNews(c.Request.Context(), &req, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Change News")
}

func (h *Handler) DisplayNews(c *gin.Context) {
	types := c.Query("type")
	n, err := h.Service.ViewNews(c.Request.Context(), types)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var news []NewsResponse

	for _, res := range *n {
		newsResponse := &NewsResponse{
			Id:          res.Id,
			Title:       res.Title,
			Content:     res.Content,
			Description: res.Description,
			Image:       res.Image,
			Category:    res.Category,
			CreatedAt:   res.CreatedAt,
		}

		news = append(news, *newsResponse)
	}

	c.JSON(http.StatusOK, news)
}

func (h *Handler) DisplayNewsById(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	res, err := h.Service.ViewNewsByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DisplayNewsCount(c *gin.Context) {
	res, err := h.Service.ViewAllNewsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
