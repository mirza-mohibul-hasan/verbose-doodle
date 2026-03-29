package service

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// QRService handles QR code generation business logic.
type QRService struct{}

// NewQRService creates a new QRService instance.
func NewQRService() *QRService {
	return &QRService{}
}

// GenerateRequest contains parameters for QR generation.
type GenerateRequest struct {
	Content         string
	ContentType     string
	Size            int
	ErrorCorrection string
	ForegroundColor string
	BackgroundColor string
}

// Generate creates a QR code PNG image from the given parameters.
func (s *QRService) Generate(req GenerateRequest) ([]byte, error) {
	// Format the content based on type
	content, err := s.formatContent(req.Content, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// Parse error correction level
	ecLevel := s.parseErrorCorrection(req.ErrorCorrection)

	// Create QR code
	qr, err := qrcode.New(content, ecLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Set colors
	fg, err := s.parseHexColor(req.ForegroundColor, color.Black)
	if err != nil {
		return nil, fmt.Errorf("invalid foreground color: %w", err)
	}
	bg, err := s.parseHexColor(req.BackgroundColor, color.White)
	if err != nil {
		return nil, fmt.Errorf("invalid background color: %w", err)
	}

	qr.ForegroundColor = fg
	qr.BackgroundColor = bg

	// Set size with default
	size := req.Size
	if size == 0 {
		size = 256
	}

	// Generate PNG
	png, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PNG: %w", err)
	}

	return png, nil
}

// formatContent converts content into the proper QR format based on type.
func (s *QRService) formatContent(content, contentType string) (string, error) {
	switch contentType {
	case "text":
		return content, nil

	case "url":
		// Ensure URL has a scheme
		if !strings.HasPrefix(content, "http://") && !strings.HasPrefix(content, "https://") {
			content = "https://" + content
		}
		return content, nil

	case "email":
		if !strings.Contains(content, "@") {
			return "", fmt.Errorf("invalid email address")
		}
		return "mailto:" + content, nil

	case "phone":
		// Strip everything except digits and + sign
		cleaned := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' || r == '+' {
				return r
			}
			return -1
		}, content)
		if len(cleaned) < 3 {
			return "", fmt.Errorf("invalid phone number")
		}
		return "tel:" + cleaned, nil

	case "wifi":
		// Parse WiFi JSON data
		var wifi struct {
			SSID       string `json:"ssid"`
			Password   string `json:"password"`
			Encryption string `json:"encryption"`
			Hidden     bool   `json:"hidden"`
		}
		if err := json.Unmarshal([]byte(content), &wifi); err != nil {
			return "", fmt.Errorf("invalid WiFi data: %w", err)
		}
		if wifi.SSID == "" {
			return "", fmt.Errorf("WiFi SSID is required")
		}
		if wifi.Encryption == "" {
			wifi.Encryption = "WPA"
		}

		hidden := "false"
		if wifi.Hidden {
			hidden = "true"
		}

		// WiFi QR code format: WIFI:T:WPA;S:network;P:password;H:hidden;;
		return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;H:%s;;",
			wifi.Encryption, wifi.SSID, wifi.Password, hidden), nil

	default:
		return content, nil
	}
}

// parseErrorCorrection converts a string level to qrcode.RecoveryLevel.
func (s *QRService) parseErrorCorrection(level string) qrcode.RecoveryLevel {
	switch strings.ToUpper(level) {
	case "L":
		return qrcode.Low
	case "M":
		return qrcode.Medium
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// parseHexColor converts a hex color string (#RRGGBB or #RRGGBBAA) to color.Color.
func (s *QRService) parseHexColor(hex string, defaultColor color.Color) (color.Color, error) {
	if hex == "" {
		return defaultColor, nil
	}

	hex = strings.TrimPrefix(hex, "#")

	switch len(hex) {
	case 6:
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid red component")
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid green component")
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid blue component")
		}
		return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil

	case 8:
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid red component")
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid green component")
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid blue component")
		}
		a, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid alpha component")
		}
		return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil

	default:
		return nil, fmt.Errorf("hex color must be 6 or 8 characters")
	}
}
