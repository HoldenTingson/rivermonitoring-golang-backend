package river

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	sortBy := c.Query("sortby")

	rivers, err := h.Service.SortRiver(c.Request.Context(), sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rivers)
}

func (h *Handler) SearchHandler(c *gin.Context) {
	location := c.Query("location")

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	SensorStatusDisconnected = "Terputus"
	SensorStatusConnected    = "Terhubung"
)

// SensorConnection represents a connected sensor with its status
type SensorConnection struct {
	LastActivity time.Time
	Channel      chan *UpdateRiver
	Status       string
	Conn         *websocket.Conn
}

// Global connection management with mutex for thread safety
var (
	sensorConnections = make(map[string]*SensorConnection)
	sensorMutex       sync.RWMutex
)

// HeartbeatInterval defines how often we check for stale connections
const HeartbeatInterval = 3 * time.Second

// ConnectionTimeout defines when a connection is considered stale
const ConnectionTimeout = 5 * time.Second

// Initialize connection monitoring
func init() {
	// Start a goroutine to monitor connections
	go monitorConnections()
}

// monitorConnections periodically checks for stale connections
func monitorConnections() {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		checkStaleConnections()
	}
}

// checkStaleConnections identifies and removes stale connections
func checkStaleConnections() {
	now := time.Now()

	sensorMutex.Lock()
	defer sensorMutex.Unlock()

	for id, conn := range sensorConnections {
		if conn.Status == SensorStatusConnected && now.Sub(conn.LastActivity) > ConnectionTimeout {
			// Connection is stale, mark as disconnected
			conn.Status = SensorStatusDisconnected
			fmt.Printf("Sensor %s marked as disconnected due to inactivity\n", id)

			// Notify any services that need to know about the disconnection
			if conn.Conn != nil {
				// Close gracefully if possible
				conn.Conn.Close()
				conn.Conn = nil
			}
		}
	}
}

// UpdateSensorActivity updates the last activity timestamp for a sensor
func UpdateSensorActivity(sensorID string) {
	sensorMutex.Lock()
	defer sensorMutex.Unlock()

	if conn, exists := sensorConnections[sensorID]; exists {
		conn.LastActivity = time.Now()
		conn.Status = SensorStatusConnected
	}
}

// GetSensorStatus returns the current status of a sensor
func GetSensorStatus(sensorID string) string {
	sensorMutex.RLock()
	defer sensorMutex.RUnlock()

	if conn, exists := sensorConnections[sensorID]; exists {
		return conn.Status
	}
	return SensorStatusDisconnected
}

// Modified handler functions
func (h *Handler) ConnectSensor(c *gin.Context) {
	riverID := c.Param("id")

	sensorMutex.RLock()
	conn, exists := sensorConnections[riverID]
	sensorMutex.RUnlock()

	if exists {
		// Check if the sensor is actually connected or just registered
		if conn.Status == SensorStatusConnected {
			fmt.Println("Sensor already connected to river:", riverID)
			c.JSON(http.StatusOK, gin.H{
				"message":   "Sensor already connected",
				"sensor_id": riverID,
				"status":    conn.Status,
			})
			return
		}

		// We have a channel but no active connection
		fmt.Println("Sensor registered but not active for river:", riverID)
		c.JSON(http.StatusOK, gin.H{
			"message":   "Sensor registered but not active",
			"sensor_id": riverID,
			"status":    conn.Status,
		})
	} else {
		// Create a new sensor connection entry
		sensorMutex.Lock()
		sensorConnections[riverID] = &SensorConnection{
			Channel:      make(chan *UpdateRiver),
			Status:       SensorStatusDisconnected,
			LastActivity: time.Now(),
		}
		dataChan := sensorConnections[riverID].Channel
		sensorMutex.Unlock()

		fmt.Println("Sensor registered for river:", riverID)

		go h.Service.UpdateRiver(context.Background(), riverID, dataChan)

		c.JSON(http.StatusOK, gin.H{
			"message":   "Sensor registered",
			"sensor_id": riverID,
			"status":    SensorStatusDisconnected,
		})
	}
}

func (h *Handler) ReceiveSensorData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensorID := c.Param("id")
	if sensorID == "" {
		conn.Close()
		c.JSON(http.StatusBadRequest, gin.H{"error": "sensor_id is required"})
		return
	}

	// Set up ping handler to detect disconnections
	conn.SetPingHandler(func(string) error {
		UpdateSensorActivity(sensorID)
		conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
		return nil
	})

	// Create or update sensor connection
	sensorMutex.Lock()
	sensorConn, exists := sensorConnections[sensorID]
	if !exists {
		sensorConn = &SensorConnection{
			Channel:      make(chan *UpdateRiver),
			Status:       SensorStatusConnected,
			LastActivity: time.Now(),
			Conn:         conn,
		}
		sensorConnections[sensorID] = sensorConn

		// Start the river update service if not already running
		go h.Service.UpdateRiver(context.Background(), sensorID, sensorConn.Channel)
	} else {
		// Update existing connection
		sensorConn.Status = SensorStatusConnected
		sensorConn.LastActivity = time.Now()
		sensorConn.Conn = conn
	}
	sensorMutex.Unlock()

	fmt.Println("New WebSocket client connected! Sensor ID:", sensorID)

	// Set up a heartbeat to detect disconnections
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					fmt.Printf("Ping failed for sensor %s: %v\n", sensorID, err)
					return
				}
			}
		}
	}()

	// Main message loop
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Connection lost for sensor %s: %v\n", sensorID, err)

			// Mark sensor as disconnected
			sensorMutex.Lock()
			if sensorConn, ok := sensorConnections[sensorID]; ok {
				sensorConn.Status = SensorStatusDisconnected
				sensorConn.Conn = nil
			}
			sensorMutex.Unlock()

			break
		}

		fmt.Println("Received from ESP32:", string(msg))
		UpdateSensorActivity(sensorID)

		height, err := strconv.ParseFloat(string(msg), 64)
		if err != nil {
			fmt.Println("Invalid height data:", err)
			continue
		}

		update := &UpdateRiver{
			Id:     sensorID,
			Height: height,
		}

		// Send update through channel
		sensorMutex.RLock()
		if sensorConn, ok := sensorConnections[sensorID]; ok {
			select {
			case sensorConn.Channel <- update:
				// Update sent successfully
			default:
				// Channel is full or blocked, log and continue
				fmt.Printf("Warning: Channel for sensor %s is blocked\n", sensorID)
			}
		}
		sensorMutex.RUnlock()
	}
}

// Add this new endpoint to check sensor status
func (h *Handler) GetSensorStatus(c *gin.Context) {
	sensorID := c.Param("id")

	status := GetSensorStatus(sensorID)

	c.JSON(http.StatusOK, gin.H{
		"sensor_id": sensorID,
		"status":    status,
	})
}

func (h *Handler) SendSensorData(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	clients[conn] = true
}
