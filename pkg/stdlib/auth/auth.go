package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/ssentinull/dealls-dating-service/internal/business/model"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/httpclient"
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/parser"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	// JWTStaticClaim const
	JWTStaticClaimID    string = "1"
	JWTStaticClaimEmail string = "admin@mail.com"
	JWTStaticClaimName  string = "admin"

	// Invalid JWT Token Message
	MessageErrorJWTInvalidToken string = "invalid jwt token"

	// ContextKey const
	ContextKeyURLPath       string = "url.path"
	ContextKeyRequestMethod string = "request.method"
	ContextKeyClaims        string = "claims"
	ContextKeyUserID        string = "user_id"
	ContextKeyUserEmail     string = "user_email"
	ContextKeyUserName      string = "user_name"
	ContextKeyJwtToken      string = "jwt_token"
	ContextXSignature       string = "X-Signature"
)

type Auth interface {
	GenerateJWTToken(user model.UserModel) (string, error)
	ParseJWTToken(tokenstring string) (*JWTClaims, error)
	ExtractJWTClaims(tokenstring string) (*JWTClaims, error)
	ParseJWTTokenWithoutExpirationCheck(tokenstring string) (*JWTClaims, error)
	GetEmail(c *gin.Context) string
	GetName(c *gin.Context) string
	GetStaticToken() string
	GetJwtToken(c *gin.Context) string
	GetSignature(c *gin.Context) string
}

type auth struct {
	opt        Options
	httpClient httpclient.HTTPClient
	json       parser.JSONParser
}

type Options struct {
	SecretKey      string `validate:"required"`
	StaticToken    string `validate:"required"`
	ExpiryDuration string `validate:"required"`
}

// JWTClaims jwt claims struct
type JWTClaims struct {
	jwt.StandardClaims
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func Init(
	opt Options,
	httpClient httpclient.HTTPClient,
	json parser.JSONParser,
) Auth {
	return &auth{
		opt:        opt,
		httpClient: httpClient,
		json:       json,
	}
}

func (a *auth) GenerateJWTToken(user model.UserModel) (string, error) {
	var expiryDuration time.Duration
	expiryDuration, err := time.ParseDuration(a.opt.ExpiryDuration)
	if err != nil {
		expiryDuration = 1 * time.Hour
	}

	claims := JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiryDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.opt.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *auth) ParseJWTToken(tokenString string) (*JWTClaims, error) {
	secretKey := a.opt.SecretKey
	staticToken := a.opt.StaticToken

	if strings.Compare(tokenString, staticToken) == 0 {
		return &JWTClaims{
			Email: JWTStaticClaimEmail,
			Name:  JWTStaticClaimName,
		}, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(MessageErrorJWTInvalidToken)
}

func (a *auth) ExtractJWTClaims(tokenString string) (*JWTClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, errors.New(MessageErrorJWTInvalidToken)
}

func (a *auth) ParseJWTTokenWithoutExpirationCheck(tokenString string) (*JWTClaims, error) {
	// check if token is valid
	claims, err := a.ParseJWTToken(tokenString)
	if err != nil {
		// check if token is invalid
		if !strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New(MessageErrorJWTInvalidToken)
		}

		// token is valid but expired, just extract the claims
		claims, err = a.ExtractJWTClaims(tokenString)
	}

	return claims, err
}

func (a *auth) GetID(c *gin.Context) string {
	return c.GetString(ContextKeyUserID)
}

func (a *auth) GetName(c *gin.Context) string {
	return c.GetString(ContextKeyUserName)
}

func (a *auth) GetEmail(c *gin.Context) string {
	return c.GetString(ContextKeyUserEmail)
}

func (a *auth) GetStaticToken() string {
	return a.opt.StaticToken
}

func (a *auth) GetJwtToken(c *gin.Context) string {
	return c.GetString(ContextKeyJwtToken)
}

func (a *auth) GetSignature(c *gin.Context) string {
	return c.GetString(ContextXSignature)
}
