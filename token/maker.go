package token

var secretKey = []byte("123456abc")

type Maker interface {
	CreateToken() (string, *Payload, error)
	VerifyToken(tokenString string) (*Payload, error)
}
