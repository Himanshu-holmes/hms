package authorization

import (
	"errors"
	"time"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
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
    ID pgtype.UUID `json:"id"`
	Username string `json:"username"`
	ExpirationDate time.Time `json:"expirationDate"`
	Type TokenType `json:"type"`
	Role model.UserRole `json:"role"`
}