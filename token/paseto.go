package token

import (
	"fmt"
	"github.com/shijting/kit/option"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	userID       int64
	username     string
	expiredAt    time.Time
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string, opts ...option.Option[PasetoMaker]) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	option.Options[PasetoMaker](opts).Apply(maker)
	return maker, nil
}

func WithPasstoUserID(userID int64) option.Option[PasetoMaker] {
	return func(t *PasetoMaker) {
		t.userID = userID
	}
}

func WithPasstoExpired(d time.Duration) option.Option[PasetoMaker] {
	return func(t *PasetoMaker) {
		t.expiredAt = time.Now().Add(d)
	}
}

func WithPasstoUsername(username string) option.Option[PasetoMaker] {
	return func(t *PasetoMaker) {
		t.username = username
	}
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken() (string, *Payload, error) {
	payload, err := NewPayload(maker.userID, maker.username, maker.expiredAt)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
