package captcha

import (
	"context"
	"time"
)

type ValidatedData interface {
	ChallengeTS() time.Time
	Hostname() string
	Action() string
}

type Validator interface {
	Validate(ctx context.Context, token string, remoteip string) (ValidatedData, error)
}
