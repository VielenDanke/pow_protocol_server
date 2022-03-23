package user

type User struct {
	password string
}

func NewUser(password string) *User {
	return &User{password: password}
}

func (u User) GetPassword() string {
	return u.password
}
