package authorization

import (
	"errors"
	"time"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrAuthorizationFailed = errors.New("authorization failed")
	ErrBadClaim = errors.New("bad jwt claim")
	ErrTokenExpired = errors.New("token is expired")
	ErrInvalidSignature = errors.New("signature is invalid")
)

type TokenType string 

var (
	AccessToken TokenType = "access-token"
	RefreshToken TokenType = "refresh-token"
)
type Info struct {
	Subject string `json:"subject"`
	ExpirationData time.Time `json:"expirationDate"`
	Type TokenType `json:"type"`
}