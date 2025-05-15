package authorization

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var secret = []byte("test")
var accessExpiration = time.Minute
var refreshExpiration = time.Minute

func TestAuth_TokenizeAndAuthorize(t *testing.T) {
	a := NewAuthorization(secret, accessExpiration, refreshExpiration)
	access, refresh, err := a.Tokenize(context.Background(), "test-subject")
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	// check if tokens are authorized
	info, err := a.Authorize(context.Background(), access)
	assert.NoError(t, err)
	assert.Equal(t, info.Subject, "test-subject")
	assert.Equal(t, info.Type, AccessToken)

	// check if the refresh token has the correct values
	info, err = a.Authorize(context.Background(), refresh)
	assert.NoError(t, err)
	assert.Equal(t, info.Subject, "test-subject")
	assert.Equal(t, info.Type, RefreshToken)
}

func TestAuth_AuthorizeBadTokenErrors(t *testing.T) {
	a := NewAuthorization(secret, accessExpiration, refreshExpiration)
	// this token is signed with another signature
	_, err := a.Authorize(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIzLTAzLTIyVDE4OjMxOjM2KzA0OjMwIiwic3ViIjoidGVzdC1zdWJqZWN0In0.qVHi1aAyRkk2SMoffcEtIesz5udHtjXqt8dq1yKUPqo")
	assert.ErrorIs(t, err, ErrInvalidSignature)
	// this token has an expired unix timestamp as its expiration date
	_, err = a.Authorize(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzk0OTUzMjMsInN1YiI6InRlc3Qtc3ViamVjdCJ9.R8QJEl63uUfJQ8vvt9p297dElaMiEH2OOVmhzySg7qA")
	assert.ErrorIs(t, err, ErrTokenExpired)
	// this token has a wrong format as it expiration date
	_, err = a.Authorize(context.Background(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiJhc2Rhc2QiLCJzdWIiOiJ0ZXN0LXN1YmplY3QifQ.xPLOY_D0x6pdDh56SB2cS5VjkN0zYnPUYgO99kSgwyU")
	assert.ErrorIs(t, err, ErrBadClaim)
}

func TestAuth_TokenizeAndRefresh(t *testing.T){
	a := NewAuthorization(secret,accessExpiration,refreshExpiration)
	access,refresh,err := a.Tokenize(context.Background(),"test-subject")
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	newToken,err := a.Refresh(context.Background(),refresh)
	assert.NoError(t,err)

	// check if the new token can be authorized
	info,err := a.Authorize(context.Background(),newToken)
	assert.NoError(t, err)
	assert.Equal(t, info.Subject, "test-subject")
	assert.Equal(t, info.Type, AccessToken)
	
}