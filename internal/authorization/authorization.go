package authorization

import "context"

type Authorization interface {
	Authorize(ctx context.Context, accessToken string) (Info, error)
	// Tokenize return access and refresh token in order
	Tokenize(ctx context.Context, subject string) (string, string, error)
	// Refresh gets the refresh token and generates a new access token
	Refresh(ctx context.Context, refreshToken string) (string, error)
}