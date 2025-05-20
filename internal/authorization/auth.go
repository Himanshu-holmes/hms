package authorization

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type auth struct {
	secret               []byte
	accessTokenDuration  time.Duration
	refreshTOkenDuration time.Duration
}

func (a *auth) Tokenize(ctx context.Context, id ,username,role string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"username": username,
		"role": role,
		"exp":  time.Now().Add(a.accessTokenDuration).Unix(),
		"type": AccessToken,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"username": username,
		"role": role,
		"exp":  time.Now().Add(a.refreshTOkenDuration).Unix(),
		"type": RefreshToken,
	})
	signedAccessToken, err := accessToken.SignedString(a.secret)
	if err != nil {
		return "", "", fmt.Errorf("error while signing access token: %w", err)
	}
	signedRefreshToken, err := refreshToken.SignedString(a.secret)
	if err != nil {
		return "", "", fmt.Errorf("error while signing refresh token: %w", err)
	}
	return signedAccessToken, signedRefreshToken, nil
}

func (a *auth)Authorize(ctx context.Context,accessToken string)(Info,error){
	token,err := a.parseAndValidate(accessToken)
	if err != nil {
		return Info{},mapToAuthErrors(err)
	}
	return claimsToInfo(token.Claims.(jwt.MapClaims))
}

func (a *auth)parseAndValidate(tokenString string)(*jwt.Token,error){
	token,err := jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) {
		if _,ok := t.Method.(*jwt.SigningMethodHMAC);!ok {
			return nil,fmt.Errorf("wrong signing algorithm: %w",ErrAuthorizationFailed)
		}
		return a.secret,nil
	})
	if err != nil {
		return nil,err
	}
	if err := a.validateToken(token);err!= nil {
		return nil,err
	}
	return token,nil
}
func (a *auth)validateToken(token *jwt.Token)error{
	if !token.Valid {
		return fmt.Errorf("token is not valid: %w",ErrAuthorizationFailed)
	}
	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrBadClaim
	}
	_,err := claims.GetSubject()
	if err != nil {
		return ErrBadClaim
	}
	return nil
}

func (a *auth) Refresh(ctx context.Context, refreshToken string) (string, error) {
	token, err := a.parseAndValidate(refreshToken)
	if err != nil {
	   return "", mapToAuthErrors(err)
	}
	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok {
	   return "", ErrBadClaim
	}
	subject, err := claims.GetSubject()
	if err != nil {
	   return "", mapToAuthErrors(err)
	}
	id, err := getUserIDFromString(subject)
	if err != nil {
	   return "", ErrBadClaim
	}
	// Get username
	username, ok := claims["username"].(string)
	if !ok {
	   return "", ErrBadClaim
	}
	roleStr, ok := claims["role"].(string)
	if !ok {
	   return "", ErrBadClaim
	}
	role := model.UserRole(roleStr)

	
 
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	   "sub":  id,
	   "username": username,
	   "role": role,
	   "exp":  time.Now().Add(a.accessTokenDuration).Unix(),
	   "type": AccessToken,
	})
	return newToken.SignedString(a.secret)
 }
func claimsToInfo(claims jwt.MapClaims)(Info,error){
	subject, err := claims.GetSubject()
	if err != nil {
		return Info{},ErrBadClaim
	}
	id, err := getUserIDFromString(subject)
	if err != nil {
		return Info{}, ErrBadClaim
	}

	// Get username
	username, ok := claims["username"].(string)
	if !ok {
		return Info{}, ErrBadClaim
	}

	expiration,err := claims.GetExpirationTime()
	if err != nil {
		return Info{},ErrBadClaim
	}
	tokenType,is := claims["type"].(string)
	if !is {
		return Info{},ErrBadClaim
	}
	roleStr, ok := claims["role"].(string)
	if !ok {
		return Info{}, ErrBadClaim
	}
	role := model.UserRole(roleStr) 
	return Info{
		ID: id,
		Username: username,
		ExpirationDate: expiration.Time,
		Type: TokenType(tokenType),
		Role: role,
	},nil
}

func mapToAuthErrors(err error)error{
	if errors.Is(err,jwt.ErrTokenExpired){
		return ErrTokenExpired
	}
	if errors.Is(err,jwt.ErrTokenSignatureInvalid){
		return ErrInvalidSignature
	}
	return ErrBadClaim
}

func NewAuthorization(secret []byte,accessTokenDuration time.Duration,refreshTOkenDuration time.Duration)Authorization{
	return &auth{
		secret:               secret,
		accessTokenDuration:  accessTokenDuration,
		refreshTOkenDuration: refreshTOkenDuration,
	}
}

func AuthWithOutDuration(secret []byte)Authorization{
	return &auth{
		secret: secret,
	}
}

func getUserIDFromString(subject string) (pgtype.UUID, error) {
    var id pgtype.UUID
	err := id.Scan(subject)
    if err != nil {
        return pgtype.UUID{}, fmt.Errorf("failed to parse UUID from string: %w", err)
    }
    return id, nil
}