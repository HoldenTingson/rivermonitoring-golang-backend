package river

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) AddRiver(c *gin.Context) {
	var req CreateRiverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.AddRiver(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Successfully Add River")
}

func (h *Handler) RemoveRiver(c *gin.Context) {
	stringId := c.Param("id")

	err := h.Service.RemoveRiver(c.Request.Context(), stringId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Remove River")
}

func (h *Handler) ChangeRiver(c *gin.Context) {
	var req UpdateRiverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringId := c.Param("id")

	err := h.Service.ChangeRiverDetail(c.Request.Context(), &req, stringId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Successfully Change River")
}

func (h *Handler) UpdateRiver(c *gin.Context) {
	h.Service.UpdateRiver(c.Request.Context())
}

func (h *Handler) DisplayRiver(c *gin.Context) {
	r, err := h.Service.ViewRiver(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []RiverResponse

	for _, res := range *r {
		riverRes := &RiverResponse{
			Id:        res.Id,
			Latitude:  res.Latitude,
			Longitude: res.Longitude,
			Location:  res.Location,
			Status:    res.Status,
			Height:    res.Height,
		}

		response = append(response, *riverRes)

	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DisplayRiverById(c *gin.Context) {
	stringId := c.Param("id")

	res, err := h.Service.ViewRiverById(c.Request.Context(), stringId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) DisplayRiverByStatus(c *gin.Context) {
	status := c.Query("status")
	r, err := h.Service.ViewRiverByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []RiverResponse

	for _, res := range *r {
		riverRes := &RiverResponse{
			Id:        res.Id,
			Latitude:  res.Latitude,
			Longitude: res.Longitude,
			Location:  res.Location,
			Status:    res.Status,
			Height:    res.Height,
		}

		response = append(response, *riverRes)

	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) SortRiversHandler(c *gin.Context) {
	sortBy := c.Query("sortby") // Get the sort field from query parameter

	rivers, err := h.Service.SortRiver(c.Request.Context(), sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rivers)
}

func (h *Handler) SearchHandler(c *gin.Context) {
	location := c.Query("location") // Get the sort field from query parameter

	rivers, err := h.Service.SearchRiver(c.Request.Context(), location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rivers)
}

func (h *Handler) DisplayRiverCount(c *gin.Context) {
	res, err := h.Service.ViewAllRiverCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan UpdateRiver)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clients[conn] = true
	defer delete(clients, conn)

	for {
		// Read updates from broadcast channel
		msg := <-broadcast

		// Send update to all connected clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				delete(clients, client)
			}
		}
	}
}
