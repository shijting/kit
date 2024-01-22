package token

import (
	"errors"
	"github.com/shijting/kit/option"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtGenerator struct {
	userID    int64
	username  string
	expiredAt time.Time
	secretKey []byte
}

// NewJWTGenerator returns a new JwtGenerator.
func NewJWTGenerator(opts ...option.Option[JwtGenerator]) Maker {
	obj := &JwtGenerator{
		//secretKey: []byte(secretKey),
		//expiredAt: time.Now().Add(viper.Get("token.access-token-duration").(time.Duration)),
	}
	option.Options[JwtGenerator](opts).Apply(obj)
	return obj
}

func WithJWTUserID(userID int64) option.Option[JwtGenerator] {
	return func(t *JwtGenerator) {
		t.userID = userID
	}
}

func WithJWTExpired(d time.Duration) option.Option[JwtGenerator] {
	return func(t *JwtGenerator) {
		t.expiredAt = time.Now().Add(d)
	}
}

func WithJWTUsername(username string) option.Option[JwtGenerator] {
	return func(t *JwtGenerator) {
		t.username = username
	}
}

func WithJWTSecretKey(secretKey string) option.Option[JwtGenerator] {
	return func(t *JwtGenerator) {
		t.secretKey = []byte(secretKey)
	}
}

// GenerateToken creates a token with payload.
func (t *JwtGenerator) CreateToken() (string, *Payload, error) {
	payload, err := NewPayload(t.userID, t.username, t.expiredAt)
	if err != nil {
		return "", nil, err
	}

	// jwt.SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", nil, err
	}
	return tokenString, payload, nil
}

// VerifyToken verifies a token is valid or invalid and returns payload if valid.
func (t *JwtGenerator) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrExpiredToken
		}
		return t.secretKey, nil
	})

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
