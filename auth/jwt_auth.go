package auth

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/wolfeidau/authinator/models"
)

var (
	// ErrTokenExpired returned when the jwt token has expired
	ErrTokenExpired = errors.New("JWT token has expired")
)

// Certs used by JWT to sign and verify tokens
type Certs struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// GenerateClaim generate a JWT token containing a claim using the supplied
// certificates and user
func GenerateClaim(certs *Certs, usr *models.User) (string, error) {
	// generate a token
	var claims = jws.Claims{
		"user_id": models.StringValue(usr.ID),
		"login":   models.StringValue(usr.Login),
		"email":   models.StringValue(usr.Email),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	j := jws.NewJWT(claims, crypto.SigningMethodRS512)
	b, err := j.Serialize(certs.PrivateKey)

	return string(b), err
}

// ValidateClaim validate the JWT token and return the user model
// decoded from the claim
func ValidateClaim(certs *Certs, token string) (*models.User, error) {

	usr := new(models.User)

	w, err := jws.ParseJWT([]byte(token))
	if err != nil {
		return nil, err
	}

	if err := w.Validate(certs.PublicKey, crypto.SigningMethodRS512); err != nil {
		return nil, err
	}

	_, isExpired := w.Claims().Expiration()

	if !isExpired {
		return nil, ErrTokenExpired
	}

	usr.Email = extractKey("email", w.Claims())
	usr.Login = extractKey("login", w.Claims())
	usr.ID = extractKey("user_id", w.Claims())

	return usr, nil
}

func extractKey(key string, claims jwt.Claims) *string {
	if claims.Has(key) {
		val := claims.Get(key)
		if s, ok := val.(string); ok {
			return models.String(s)
		}
	}
	return nil
}
