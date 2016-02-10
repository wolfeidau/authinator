package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/wolfeidau/authinator/auth"
	"github.com/wolfeidau/authinator/store/users"
)

var userHash = "LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"

func TestAuthenticateUser(t *testing.T) {

	certs, err := generateTestCerts()
	if err != nil {
		t.Errorf("error generating test certs %v", err)
	}

	store := users.NewUserStoreLocal()

	store.Create(NewUser())

	ws := NewAuthResource(store, nil, certs)

	req := newFormRequest("POST", "http://api.his.com/users", bytes.NewBufferString("login=wolfeidau&password=Somewh3r3 there is a cow!"))

	recorder, resp := newResponse()

	ws.authenticateUser(req, resp)

	if recorder.Code != 200 {
		t.Errorf("expected 200 got %d %s", recorder.Code, recorder.Body.String())
	}

	if recorder.Header().Get("Authorization") == "" {
		t.Errorf("expected authorization header to exist")
	}

}

func generateTestCerts() (*auth.Certs, error) {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.DNSNames = append(template.DNSNames, "example.com")

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)

	if err != nil {
		return nil, err
	}

	var cbuf bytes.Buffer
	pem.Encode(&cbuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	var kbuf bytes.Buffer
	pem.Encode(&kbuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	privateKey, err := crypto.ParseRSAPrivateKeyFromPEM(kbuf.Bytes())
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.ParseRSAPublicKeyFromPEM(cbuf.Bytes())
	if err != nil {
		return nil, err
	}

	return &auth.Certs{PrivateKey: privateKey, PublicKey: publicKey}, nil
}
