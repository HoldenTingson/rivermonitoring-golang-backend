package carrousel

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

func (h *Handler) AddCarrousel(c *gin.Context) {
	var r CreateCarrouselRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.AddCarrousel(c.Request.Context(), &r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) RemoveCarrousel(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.RemoveCarrousel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Remove Carrousel")
}

func (h *Handler) ChangeCarrousel(c *gin.Context) {
	var req UpdateCarrouselRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	err := h.Service.ChangeCarrousel(c.Request.Context(), &req, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Change Carrousel")
}

func (h *Handler) DisplayCarrousel(c *gin.Context) {
	n, err := h.Service.ViewCarrousel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var Carrousel []CarrouselResponse

	for _, res := range *n {
		CarrouselResponse := &CarrouselResponse{
			Id:    res.Id,
			Title: res.Title,
			Desc:  res.Desc,
			Image: res.Image,
			Date:  res.Date,
		}

		Carrousel = append(Carrousel, *CarrouselResponse)
	}

	c.JSON(http.StatusOK, Carrousel)
}

func (h *Handler) DisplayCarrouselById(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	res, err := h.Service.ViewCarrouselByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DisplayCarrouselByIdAdmin(c *gin.Context) {
	stringId := c.Param("id")
	id, _ := strconv.Atoi(stringId)

	res, err := h.Service.ViewCarrouselByIDAdmin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DisplayCarrouselCount(c *gin.Context) {
	res, err := h.Service.ViewAllCarrouselCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
