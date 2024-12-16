package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMake is a PASETO token maker
type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {
	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for specific username, role and duration
func (maker *PasetoMaker) CreateToken(username, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symetricKey, payload, nil)

	return token, payload, err
}

// VerifyToken checks if the input token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	if err := maker.paseto.Decrypt(token, maker.symetricKey, payload, nil); err != nil {
		return nil, ErrInvalidToken
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
