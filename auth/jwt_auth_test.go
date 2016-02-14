package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wolfeidau/authinator/models"
)

func TestGenerateAndValidateClaim(t *testing.T) {

	certs, err := GenerateTestCerts()
	if assert.NoError(t, err) {

		claim, err := GenerateClaim(certs, &models.User{
			ID:       models.String("123"),
			Login:    models.String("wolfeidau"),
			Email:    models.String("mark@wolfe.id.au"),
			Name:     models.String("Mark Wolfe"),
			Password: models.String("LkSquwzxdgzSTqqc7Rku5NF8/uR7TBFO1IRF1Yj2c0sM4HEVGgp0bJadWtRAaINP"), //Somewh3r3 there is a cow!
		})

		if assert.NoError(t, err) {

			usr, err := ValidateClaim(certs, claim)
			if assert.NoError(t, err) {
				assert.NotNil(t, usr)
				assert.Equal(t, "123", models.StringValue(usr.ID))
				assert.Equal(t, "wolfeidau", models.StringValue(usr.Login))
				assert.Equal(t, "mark@wolfe.id.au", models.StringValue(usr.Email))
			}
		}
	}
}
