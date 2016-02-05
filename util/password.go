package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

var (
	// N is a CPU/memory cost
	N int32 = 16384

	R int32 = 8
	// P Parallelization
	P int32 = 1
)

// HashPassword returns a password has created using scrypt.
// return [salt] + scrypt([salt], [credential], N=16384, r=8, p=1, keyLen=32);
func HashPassword(password string) (string, error) {
	var hash string

	salt, err := generateSalt()
	if err != nil {
		return hash, err
	}
	dk, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return hash, err
	}

	rhash := append(salt, dk...)

	hash = base64.StdEncoding.EncodeToString(rhash)

	return hash, nil
}

// CompareHashPassword compares password hash by decoding and extracting the seed
// then calculating the hash and using constant time comparison to compare then.
func CompareHashPassword(password, hash string) (bool, error) {

	rhash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	salt := rhash[:16]
	odk := rhash[16:]

	dk, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare(dk, odk) == 1 {
		return true, nil
	}

	return false, nil
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}
