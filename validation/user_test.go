package validation

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/validation/field"
)

func TestValidateUpdateUser(t *testing.T) {
	testCases := []struct {
		newUser  *models.User
		oldUser  *models.User
		expected field.ErrorList
	}{
		{
			newUser:  models.NewUser("123", "wolfeidau", "mark@wolfe.id.au", "Mark Wolfe"),
			oldUser:  models.NewUser("123", "wolfeidau", "mark@wolfe.id.au", "Mark Wolfe"),
			expected: field.ErrorList{},
		},
		{
			newUser: &models.User{
				Name: models.String("Mark Wolf"),
			},
			oldUser:  models.NewUser("123", "wolfeidau", "mark@wolfe.id.au", "Mark Wolfe"),
			expected: field.ErrorList{},
		},
		{
			newUser: &models.User{
				Name:     models.String("Mark Wolf"),
				Password: models.String("Somewh3r3 there is a cow!"),
			},
			oldUser: models.NewUser("123", "wolfeidau", "mark@wolfe.id.au", "Mark Wolfe"),
			expected: field.ErrorList{
				&field.Error{Type: field.ErrorTypeForbidden, Field: field.NewPath("User", "Password").String(), BadValue: "", Detail: "User updates must not supply Password"},
			},
		},
	}

	for _, testCase := range testCases {
		errList := ValidateUserUpdate(testCase.newUser, testCase.oldUser)

		if !reflect.DeepEqual(errList, testCase.expected) {
			t.Errorf("expected\n%s\ngot\n%s\n", toJSON(testCase.expected), toJSON(errList))
		}
	}
}

func TestValidateUserRegister(t *testing.T) {
	testCases := []struct {
		newUser  *models.User
		expected field.ErrorList
	}{
		{
			newUser: &models.User{
				Login:    models.String("wolfeidau"),
				Email:    models.String("mark@wolfe.id.au"),
				Name:     models.String("Mark Wolf"),
				Password: models.String("Somewh3r3 there is a cow!"),
			},
			expected: field.ErrorList{},
		},
		{
			newUser: &models.User{
				ID:       models.String("123"),
				Login:    models.String("wolfeidau"),
				Email:    models.String("mark@wolfe.id.au"),
				Name:     models.String("Mark Wolf"),
				Password: models.String("Somewh3r3 there is a cow!"),
			},
			expected: field.ErrorList{
				&field.Error{Type: field.ErrorTypeForbidden, Field: field.NewPath("User", "ID").String(), BadValue: "", Detail: "User updates must not supply ID"},
			},
		},
		{
			newUser: &models.User{
				ID:    models.String("123"),
				Name:  models.String("Mark Wolf"),
				Email: models.String("mark@wolfe.id.au"),
			},
			expected: field.ErrorList{
				&field.Error{Type: field.ErrorTypeForbidden, Field: field.NewPath("User", "ID").String(), BadValue: "", Detail: "User updates must not supply ID"},
				&field.Error{Type: field.ErrorTypeRequired, Field: field.NewPath("User", "Login").String(), BadValue: "", Detail: "User updates must supply Login"},
				&field.Error{Type: field.ErrorTypeRequired, Field: field.NewPath("User", "Password").String(), BadValue: "", Detail: "User updates must supply Password"},
				&field.Error{Type: field.ErrorTypeInvalid, Field: field.NewPath("User", "Login").String(), BadValue: "", Detail: "User: Login must be between 5 and 255 characters"},
				&field.Error{Type: field.ErrorTypeInvalid, Field: field.NewPath("User", "Password").String(), BadValue: "", Detail: "User: Password must be between 5 and 255 characters"},
			},
		},
	}

	for _, testCase := range testCases {
		errList := ValidateUserRegister(testCase.newUser)

		if !reflect.DeepEqual(errList, testCase.expected) {
			t.Errorf("expected\n%s\ngot\n%s\n", toJSON(testCase.expected), toJSON(errList))
		}
	}
}

func toJSON(o interface{}) string {
	buf, err := json.Marshal(o)

	if err != nil {
		return ""
	}

	return string(buf)
}
