package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gos/app/repo"
)

// TokenHeader token for auth
const TokenHeader = "x-access-token"

// AccessTokenExpirationMinutes is the expiry time for the token
const AccessTokenExpirationMinutes = 300

// IAuth is an interface for handling auth
type IAuth interface {
	AuthenticateUser(ctx *gin.Context, accessToken string) (string, error)
	GetJWTKey() []byte
}

// Claims is used for auth
type Claims struct {
	UserId int64  `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

// Auth keeps the auth and a secret key
type Auth struct {
	jwtKey []byte
	repo   repo.IAppRepo
}

// NewAuth create a new auth instance
func NewAuth(repo repo.IAppRepo, jwtKey string) *Auth {
	return &Auth{
		repo:   repo,
		jwtKey: []byte(jwtKey),
	}
}

// GetJWTKey returns the current key
func (auth *Auth) GetJWTKey() []byte {
	return auth.jwtKey
}

// AuthenticateUser will auth user and returns a access token with expiry
func (auth *Auth) AuthenticateUser(ctx *gin.Context, accessToken string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.jwtKey, nil
	})

	if err != nil {
		return "", errors.Wrap(err, "access token is invalid")
	}

	if !tkn.Valid {
		return "", errors.New("access token is invalid")
	}

	claim := tkn.Claims.(*Claims)
	ctx.Set("claims", claim)

	return accessToken, nil
}

var _ = (*IAuth)(nil)
