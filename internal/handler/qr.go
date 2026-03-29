package handler

import (
	"fmt"
	"net/http"
	"time"

	"qrcode-generator/internal/model"
	"qrcode-generator/internal/service"

	"github.com/gin-gonic/gin"
)

// QRHandler handles HTTP requests for QR code operations.
type QRHandler struct {
	service *service.QRService
}

// NewQRHandler creates a new QRHandler with the given service.
func NewQRHandler(svc *service.QRService) *QRHandler {
	return &QRHandler{service: svc}
}

// Generate handles POST /api/generate — generates a QR code image.
func (h *QRHandler) Generate(c *gin.Context) {
	var req model.QRRequest

	// Bind and validate JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation_error",
			Message: fmt.Sprintf("Invalid request: %s", err.Error()),
		})
		return
	}

	// Apply defaults
	if req.Size == 0 {
		req.Size = 256
	}
	if req.ErrorCorrection == "" {
		req.ErrorCorrection = "M"
	}

	// Generate QR code
	png, err := h.service.Generate(service.GenerateRequest{
		Content:         req.Content,
		ContentType:     req.ContentType,
		Size:            req.Size,
		ErrorCorrection: req.ErrorCorrection,
		ForegroundColor: req.ForegroundColor,
		BackgroundColor: req.BackgroundColor,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "generation_error",
			Message: fmt.Sprintf("Failed to generate QR code: %s", err.Error()),
		})
		return
	}

	// Set response headers for PNG image
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=\"qrcode.png\"")
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "image/png", png)
}

// Health handles GET /api/health — returns server health status.
func (h *QRHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, model.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
