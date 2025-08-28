package cfturnstile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/International-Combat-Archery-Alliance/captcha"
	"github.com/google/uuid"
)

const (
	turnstileURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

var _ captcha.Validator = &Validator{}

type HTTPDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type Validator struct {
	client HTTPDoer
}

func NewValidator(client HTTPDoer) *Validator {
	return &Validator{
		client: client,
	}
}

func (v *Validator) Validate(ctx context.Context, token string, remoteip string) (captcha.ValidatedData, error) {
	idempotencyKey := uuid.New()

	form := url.Values{
		"secret":          {"TODO"},
		"response":        {token},
		"idempotency_key": {idempotencyKey.String()},
	}

	if remoteip != "" {
		form.Add("remoteip", remoteip)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to make http request: %w", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}

	var cfResp CFResponse
	err = json.NewDecoder(resp.Body).Decode(&cfResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode resp: %w", err)
	}

	if !cfResp.CFSuccess || len(cfResp.CFErrorCodes) > 0 {
		return nil, fmt.Errorf("invalid cf token. error codes: %s", strings.Join(cfResp.CFErrorCodes, ", "))
	}

	return &cfResp, nil
}

var _ captcha.ValidatedData = &CFResponse{}

type CFResponse struct {
	CFSuccess     bool      `json:"success"`
	CFChallengeTs time.Time `json:"challenge_ts"`
	CFHostname    string    `json:"hostname"`
	CFErrorCodes  []string  `json:"error-codes"`
	CFAction      string    `json:"action"`
	CFData        string    `json:"cdata"`
}

func (c *CFResponse) Action() string {
	return c.CFAction
}

func (c *CFResponse) ChallengeTS() time.Time {
	return c.CFChallengeTs
}

func (c *CFResponse) Hostname() string {
	return c.CFHostname
}
