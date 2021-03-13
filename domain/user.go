package domain

import (
	"errors"
	"github.com/google/uuid"
)

type User struct {
	Id       string
	Name     string
	Login    string
	NickName string
	Points   int
	Weekend  Weekend
}

func userValidate(name string, login string, nickname string) error {
	if len(name) < 3 {
		return errors.New("name must be more 3 characters")
	}

	if len(login) < 3 {
		return errors.New("login must be more 3 characters")
	}

	if len(nickname) < 3 {
		return errors.New("nickname must be more 3 characters")
	}

	return nil
}

func NewUser(name string, login string, nickname string) (*User, error) {
	err := userValidate(name, login, nickname)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:       uuid.NewString(),
		Name:     name,
		Login:    login,
		NickName: nickname,
		Points:   0,
		Weekend:  Weekend{},
	}, nil
}

func (u *User) Edit(name string, login string, nickname string) error {
	err := userValidate(name, login, nickname)
	if err != nil {
		return err
	}

	u.Name = name
	u.Login = login
	u.NickName = nickname

	return nil
}

func (u *User) AddPoint(point int) {
	u.Points += point
}

func (u *User) DeletePoint(point int) {
	u.Points -= point
}
