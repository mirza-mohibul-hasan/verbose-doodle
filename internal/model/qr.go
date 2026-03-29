package model

// QRRequest represents the incoming request to generate a QR code.
type QRRequest struct {
	Content         string `json:"content" binding:"required,max=2048"`
	ContentType     string `json:"content_type" binding:"required,oneof=text url email phone wifi"`
	Size            int    `json:"size" binding:"omitempty,min=128,max=1024"`
	ErrorCorrection string `json:"error_correction" binding:"omitempty,oneof=L M Q H"`
	ForegroundColor string `json:"foreground_color" binding:"omitempty"`
	BackgroundColor string `json:"background_color" binding:"omitempty"`
}

// WiFiData holds WiFi-specific fields embedded in content as JSON.
type WiFiData struct {
	SSID       string `json:"ssid"`
	Password   string `json:"password"`
	Encryption string `json:"encryption"` // WPA, WEP, nopass
	Hidden     bool   `json:"hidden"`
}

// HealthResponse is the response for the health check endpoint.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse is a standardized API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
