# QR Code Generator

A production-ready QR Code generator built with **Go** and **Gin**, featuring a premium dark-mode web interface.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)

## Features

### Core

- 📝 **Multiple content types**: Text, URL, Email, Phone, WiFi
- 🎨 **Custom colors**: Foreground & background via hex color picker
- 📐 **Adjustable size**: 128px to 1024px
- 🛡️ **Error correction levels**: L, M, Q, H
- 📥 **PNG download** with timestamped filenames

### Production-Ready

- ⚡ **Rate limiting**: Token bucket (10 req/s per IP)
- 🧹 **Input sanitization**: Max 2048 chars, type validation
- 📊 **Structured logging**: JSON via `log/slog`
- 🔄 **Graceful shutdown**: SIGINT/SIGTERM handling
- 🌐 **CORS**: Configurable cross-origin support
- ❤️ **Health check**: `GET /api/health`

### Frontend

- 🌙 Dark glassmorphic design with animated background
- ✨ Micro-animations and smooth transitions
- 📱 Fully responsive (mobile-first)
- 🗂️ Tab-based content type selector
- 👁️ Live preview panel

## Quick Start

### Prerequisites

- [Go 1.21+](https://go.dev/dl/) installed

### Run

```bash
# Clone and enter directory
cd qrcode-generator

# Download dependencies
go mod tidy

# Run the server
go run main.go
```

Open **http://localhost:8080** in your browser.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |

## API Reference

### Generate QR Code

```http
POST /api/generate
Content-Type: application/json
```

**Request Body:**

```json
{
  "content": "https://example.com",
  "content_type": "url",
  "size": 512,
  "error_correction": "M",
  "foreground_color": "#ffffff",
  "background_color": "#000000"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `content` | string | ✅ | Content to encode (max 2048 chars) |
| `content_type` | string | ✅ | `text`, `url`, `email`, `phone`, `wifi` |
| `size` | int | ❌ | Image size in px (128–1024, default: 256) |
| `error_correction` | string | ❌ | `L`, `M`, `Q`, `H` (default: `M`) |
| `foreground_color` | string | ❌ | Hex color (default: `#000000`) |
| `background_color` | string | ❌ | Hex color (default: `#ffffff`) |

**WiFi Content Format** (JSON string in `content`):

```json
{
  "ssid": "MyNetwork",
  "password": "secret123",
  "encryption": "WPA",
  "hidden": false
}
```

**Response:** `image/png` (200) or JSON error (4xx/5xx)

### Health Check

```http
GET /api/health
```

```json
{
  "status": "ok",
  "timestamp": "2026-03-29T15:00:00Z"
}
```

## Project Structure

```
qrcode-generator/
├── main.go                      # Entry point, router, graceful shutdown
├── go.mod                       # Go module definition
├── internal/
│   ├── handler/
│   │   └── qr.go                # HTTP handlers
│   ├── middleware/
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Structured logging middleware
│   │   └── ratelimit.go         # Per-IP rate limiting
│   ├── model/
│   │   └── qr.go                # Request/response models
│   └── service/
│       └── qr.go                # QR generation business logic
├── static/
│   ├── index.html               # Web UI
│   ├── css/
│   │   └── style.css            # Design system
│   └── js/
│       └── app.js               # Frontend logic
└── README.md
```

## License

MIT
