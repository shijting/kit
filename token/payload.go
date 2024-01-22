package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is expired")
	ErrExpiredToken = errors.New("invalid token")
)

type Payload struct {
	TokenID   uuid.UUID `json:"token_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func NewPayload(userID int64, username string, expiredAt time.Time) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		TokenID:   tokenID,
		UserID:    userID,
		Username:  username,
		ExpiredAt: expiredAt,
		IssuedAt:  time.Now(),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
