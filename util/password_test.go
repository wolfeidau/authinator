package util

import "testing"

func TestCompareHashPassword(t *testing.T) {

	pass := "Somewh3r3 there is a cow!"

	phash, err := HashPassword(pass)

	if err != nil {
		t.Errorf("error hashing pass %v", err)
	}

	ok, err := CompareHashPassword(pass, phash)
	if err != nil {
		t.Errorf("error comparing hash password %v", err)
	}

	if !ok {
		t.Errorf("expected true got %v", ok)
	}
}
