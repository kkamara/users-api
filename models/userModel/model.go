package userModel

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/kkamara/users-api/config"
	"github.com/kkamara/users-api/schemas/userSchema"
)

func GenerateUsername(first_name, last_name string) string {
	p, _ := rand.Prime(rand.Reader, 64)
	return fmt.Sprintf(
		"%s%s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %s", first_name, last_name)))[:9],
		p.String()[1:5],
	)
}

func ValidateCreate(newUser *userSchema.UserSchema) (errors []string) {
	if len(newUser.FirstName) < 1 || len(newUser.FirstName) >= 50 {
		errors = append(errors, "The first_name field must be between 1-50 characters in length.")
	}
	if len(newUser.LastName) < 1 || len(newUser.LastName) >= 50 {
		errors = append(errors, "The last_name field must be between 1-50 characters in length.")
	}
	return
}

func Create(newUser *userSchema.UserSchema) (user *userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	res := db.Create(&newUser)
	user = newUser
	err = res.Error
	return
}
