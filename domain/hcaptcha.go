package domain

import "context"

type HCaptchaClient interface {
	Verify(ctx context.Context, clientResponse string) (success bool, err error)
}
