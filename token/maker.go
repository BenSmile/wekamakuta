package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for specific username and duration
	CreateToken(username, role string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the input token is valid or not
	VerifyToken(token string) (*Payload, error)
}
