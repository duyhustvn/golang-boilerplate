package authmodel

// User kafka message
// {"username": "duyle", "password": "123"}
type User struct {
	Username string
	Password string
}
