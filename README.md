# CAPTCHA

A Go library providing a generic interface for CAPTCHA validation with pluggable implementations.

## Features

- Generic `Validator` interface for CAPTCHA validation
- Cloudflare Turnstile implementation included
- Context-aware validation
- Structured validation data with timestamp, hostname, and action information

## Installation

```bash
go get github.com/International-Combat-Archery-Alliance/captcha
```

## Usage

### Basic Usage with Cloudflare Turnstile

```go
package main

import (
    "context"
    "net/http"

    "github.com/International-Combat-Archery-Alliance/captcha"
    "github.com/International-Combat-Archery-Alliance/captcha/cfturnstile"
)

func main() {
    // Create a Turnstile validator with HTTP client
    validator := cfturnstile.NewValidator(&http.Client{})
    
    // Validate a CAPTCHA token
    ctx := context.Background()
    token := "your-captcha-token"
    remoteIP := "192.168.1.1" // optional
    
    data, err := validator.Validate(ctx, token, remoteIP)
    if err != nil {
        // Handle validation error
        return
    }
    
    // Access validation data
    timestamp := data.ChallengeTS()
    hostname := data.Hostname()
    action := data.Action()
}
```

### Custom Implementation

Implement the `captcha.Validator` interface for other CAPTCHA providers:

```go
type CustomValidator struct {
    // your fields
}

func (v *CustomValidator) Validate(ctx context.Context, token string, remoteip string) (captcha.ValidatedData, error) {
    // your implementation
}
```

## Interface

### Validator

```go
type Validator interface {
    Validate(ctx context.Context, token string, remoteip string) (ValidatedData, error)
}
```

### ValidatedData

```go
type ValidatedData interface {
    ChallengeTS() time.Time
    Hostname() string
    Action() string
}
```

## Cloudflare Turnstile

The included Cloudflare Turnstile implementation:

- Validates tokens against Cloudflare's API
- Supports optional remote IP validation
- Uses UUID-based idempotency keys
- Returns structured response data

## Requirements

- Go 1.24.6 or later

## License

See [LICENSE](LICENSE) file.
