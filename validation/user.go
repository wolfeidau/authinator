package validation

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	"github.com/wolfeidau/authinator/models"
	"github.com/wolfeidau/authinator/validation/field"
)

// ValidateUserUpdate validate user update requests
func ValidateUserUpdate(newUser, oldUser *models.User) field.ErrorList {
	allErrs := field.ErrorList{}

	path := field.NewPath("User")

	allErrs = append(allErrs, validateImmutibleFields(newUser, oldUser, path, []string{"ID", "Email", "Login"})...)
	allErrs = append(allErrs, validateInvalidFields(newUser, path, []string{"Password"})...)

	return allErrs
}

// ValidateUserRegister validate user registration requests
func ValidateUserRegister(newUser *models.User) field.ErrorList {
	allErrs := field.ErrorList{}

	path := field.NewPath("User")

	allErrs = append(allErrs, validateInvalidFields(newUser, path, []string{"ID"})...)
	allErrs = append(allErrs, validateRequiredFields(newUser, path, []string{"Email", "Login", "Password"})...)

	allErrs = append(allErrs, validateFieldLength(newUser.Email, path, 5, 255, "Email")...)
	allErrs = append(allErrs, validateFieldLength(newUser.Login, path, 5, 255, "Login")...)
	allErrs = append(allErrs, validateFieldLength(newUser.Password, path, 5, 255, "Password")...)

	return allErrs
}

func validateImmutibleFields(new, old interface{}, fldPath *field.Path, fields []string) field.ErrorList {
	allErrs := field.ErrorList{}

	olds := structs.New(old)
	news := structs.New(new)

	for _, f := range fields {

		// ommitting immutable fields is OK
		if news.Field(f).IsZero() {
			continue
		}

		if !reflect.DeepEqual(news.Field(f).Value(), olds.Field(f).Value()) {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child(f), fmt.Sprintf("%s updates must not change %s", fldPath.String(), news.Field(f).Name())))
		}
	}

	return allErrs
}

func validateInvalidFields(new interface{}, fldPath *field.Path, fields []string) field.ErrorList {
	allErrs := field.ErrorList{}

	news := structs.New(new)

	for _, f := range fields {

		// ommitting immutable fields is OK
		if news.Field(f).IsZero() {
			continue
		}

		allErrs = append(allErrs, field.Forbidden(fldPath.Child(f), fmt.Sprintf("%s updates must not supply %s", fldPath.String(), news.Field(f).Name())))
	}

	return allErrs
}

func validateRequiredFields(new interface{}, fldPath *field.Path, fields []string) field.ErrorList {
	allErrs := field.ErrorList{}

	news := structs.New(new)

	for _, f := range fields {

		if news.Field(f).IsZero() {
			allErrs = append(allErrs, field.Required(fldPath.Child(f), fmt.Sprintf("%s updates must supply %s", fldPath.String(), news.Field(f).Name())))
		}

	}

	return allErrs
}

func validateFieldLength(value *string, fldPath *field.Path, min, max int, fieldName string) field.ErrorList {
	allErrs := field.ErrorList{}

	if v := models.StringValue(value); !checkFieldLength(v, min, max) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldName), v, fmt.Sprintf("%s: %s must be between %d and %d characters", fldPath.String(), fieldName, min, max)))
	}

	return allErrs
}

func checkFieldLength(value string, min, max int) bool {
	if len(value) < min {
		return false
	}

	if len(value) > max {
		return false
	}

	return true
}
