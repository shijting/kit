package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	testCases := []struct {
		name string

		userID     int64
		username   string
		duration   time.Duration
		inputToken string

		wantExpiredAt time.Time
		wantGenErr    error
		wantCheckErr  error
	}{
		{
			name:          "ok",
			userID:        6,
			username:      "Alice",
			duration:      time.Second * 10,
			wantExpiredAt: time.Now().Add(time.Second * 10),
		},
		{
			name:          "expired token",
			userID:        6,
			username:      "Alice",
			duration:      -time.Second * 10,
			wantExpiredAt: time.Now().Add(-time.Second * 10),
			wantCheckErr:  ErrExpiredToken,
		},
		{
			name:          "invalid token",
			userID:        6,
			username:      "Alice",
			duration:      time.Second * 10,
			inputToken:    "123456",
			wantExpiredAt: time.Now().Add(time.Second * 10),
			wantCheckErr:  ErrInvalidToken,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			secretKey := "123456"
			jwt := NewJWTGenerator(WithJWTSecretKey(secretKey), WithJWTUserID(testCase.userID), WithJWTExpired(testCase.duration), WithJWTUsername(testCase.username))
			token, payload, err := jwt.CreateToken()

			assert.Equal(t, testCase.wantGenErr, err)
			if err != nil {
				return
			}
			assert.NotEmpty(t, token)
			assert.Equal(t, testCase.userID, payload.UserID)
			assert.Equal(t, testCase.username, payload.Username)
			require.WithinDuration(t, payload.ExpiredAt, testCase.wantExpiredAt, time.Second)

			if testCase.inputToken != "" {
				token = testCase.inputToken
			}
			vp, err := jwt.VerifyToken(token)
			assert.Equal(t, testCase.wantCheckErr, err)
			if err != nil {
				return
			}
			assert.NotEmpty(t, vp)
			assert.Equal(t, vp.UserID, testCase.userID)
			assert.Equal(t, vp.Username, testCase.username)
			require.WithinDuration(t, vp.ExpiredAt, testCase.wantExpiredAt, time.Second)
		})
	}
}
