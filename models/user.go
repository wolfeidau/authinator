package models

// String returns a pointer to of the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// NewUser helper method to create a user
func NewUser(id, login, email, name string) *User {
	return &User{
		ID:    String(id),
		Login: String(login),
		Email: String(email),
		Name:  String(name),
	}
}

// User represents a authinator user.
type User struct {
	ID       *string `json:"id,omitempty"`
	Login    *string `json:"login,omitempty"`
	Email    *string `json:"email,omitempty"`
	Name     *string `json:"name,omitempty"`
	Password *string `json:"password,omitempty"`
}
