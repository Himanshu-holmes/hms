package authorization

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/stretchr/testify/assert"
)

var secret = []byte("test-secret")
var accessExpiration = time.Minute
var refreshExpiration = time.Minute

func TestAuth_TokenizeAndAuthorize(t *testing.T) {
	a := NewAuthorization(secret, accessExpiration, refreshExpiration)
	id := uuid.New().String()
	username := "testuser"
	role := string(model.RoleDoctor)

	access, refresh, err := a.Tokenize(context.Background(), id, username, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	// Validate access token
	info, err := a.Authorize(context.Background(), access)
	assert.NoError(t, err)
	assert.Equal(t, id, info.ID.String())
	assert.Equal(t, username, info.Username)
	assert.Equal(t, AccessToken, info.Type)
	assert.Equal(t, model.UserRole(role), info.Role)

	// Validate refresh token
	info, err = a.Authorize(context.Background(), refresh)
	assert.NoError(t, err)
	assert.Equal(t, id, info.ID.String())
	assert.Equal(t, username, info.Username)
	assert.Equal(t, RefreshToken, info.Type)
	assert.Equal(t, model.UserRole(role), info.Role)
}

func TestAuth_AuthorizeBadTokenErrors(t *testing.T) {
	expiredTime := time.Now().Add(-time.Hour).Minute()
	a := NewAuthorization(secret, time.Duration(expiredTime) , refreshExpiration)

	// Token signed with wrong key
	_, err := a.Authorize(context.Background(),
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30")
	assert.ErrorIs(t, err, ErrInvalidSignature)

	id := uuid.New().String()
	username := "testuser"
	role := string(model.RoleDoctor)

	access, _, err := a.Tokenize(context.Background(), id, username, role)
   assert.NoError(t, err)
	// Expired token
	_, err = a.Authorize(context.Background(),
		access)
	assert.ErrorIs(t, err, ErrTokenExpired)

	// Malformed expiration
	_, err = a.Authorize(context.Background(),
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiJhc2Rhc2QiLCJzdWIiOiIxMjMiLCJ1c2VybmFtZSI6InRlc3QiLCJyb2xlIjoiU3R1ZGVudCIsInR5cGUiOiJBY2Nlc3NUb2tlbiJ9.5Xq-yJb5B5K2cJKFuznhro-Vop5MmYkL_0ZpGe9MHE8")
	assert.ErrorIs(t, err, ErrInvalidSignature)
}

func TestAuth_TokenizeAndRefresh(t *testing.T) {
	a := NewAuthorization(secret, accessExpiration, refreshExpiration)
	id := uuid.New().String()
	username := "testuser"
	role := string(model.RoleReceptionist)

	access, refresh, err := a.Tokenize(context.Background(), id, username, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	newAccess, err := a.Refresh(context.Background(), refresh)
	assert.NoError(t, err)

	info, err := a.Authorize(context.Background(), newAccess)
	assert.NoError(t, err)
	assert.Equal(t, id, info.ID.String())
	assert.Equal(t, username, info.Username)
	assert.Equal(t, AccessToken, info.Type)
	assert.Equal(t, model.UserRole(role), info.Role)
}
