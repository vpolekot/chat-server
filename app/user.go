package app

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"UserName"`
	Password string `json:"Password"`
}

var users = make(map[string]*User)

func generateToken(id string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expDate := time.Now().Add(time.Hour * 24)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["admin"] = true
	claims["exp"] = expDate.Unix()

	t, err := token.SignedString([]byte(uuid.New().String()))
	if err != nil {
		return "", "", err
	}
	return t, expDate.UTC().String(), nil
}

func getUserId(credentials requestCredentials, users map[string]*User) (string, bool) {
	u, ok := users[credentials.UserName]
	if ok && u.Password == credentials.Password {
		return u.ID, true
	}
	return "", false
}

func validateCredentials(cred requestCredentials) error {
	validate := validator.New()
	err := validate.Struct(cred)
	if err != nil {
		return errors.New("Bad request, empty username or id")
	}
	return nil
}
