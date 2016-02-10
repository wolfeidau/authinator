package auth

import (
	"crypto/rsa"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/wolfeidau/authinator/models"
)

// Certs used by JWT to sign and verify tokens
type Certs struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// GenerateClaim generate a JWT token containing a cleam using the suplied
// certs and user
func GenerateClaim(certs *Certs, usr *models.User) (string, error) {
	// generate a token
	var claims = jws.Claims{
		"email": models.StringValue(usr.Email),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	j := jws.NewJWT(claims, crypto.SigningMethodRS512)
	b, err := j.Serialize(certs.PrivateKey)

	return string(b), err
}
