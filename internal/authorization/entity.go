package authorization

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/model"
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
    ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	ExpirationDate time.Time `json:"expirationDate"`
	Type TokenType `json:"type"`
	Role model.UserRole `json:"role"`
}