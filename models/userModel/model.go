package userModel

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/bxcodec/faker/v3"

	"github.com/kkamara/users-api/config"
	"github.com/kkamara/users-api/schemas/userSchema"
)

func Get(username string) (user *userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	res := db.Where("username = ?", username).First(&user)
	err = res.Error
	return
}

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

func ValidateUpdate(newUser *userSchema.UserSchema) (errors []string) {
	return ValidateCreate(newUser)
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

func Update(username string, updateUser *userSchema.UserSchema) (user *userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	newUser := &userSchema.UserSchema{
		FirstName: updateUser.FirstName,
		LastName:  updateUser.LastName,
		DarkMode:  updateUser.DarkMode,
	}
	var schema *userSchema.UserSchema
	res := db.Model(&schema).Select(
		"first_name", "last_name", "dark_mode",
	).Where("username = ?", username).Updates(newUser)
	user = updateUser
	err = res.Error
	return
}

func GetAll() (users []*userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	res := db.Find(&users)
	err = res.Error
	return
}

func DelUser(username string) (err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	var schema *userSchema.UserSchema
	res := db.Where("username = ?", username).Delete(&schema)
	err = res.Error
	return
}

func FindUsers(query string) (users []*userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}
	formattedQuery := fmt.Sprintf("%%%s%%", query)
	res := db.Where(
		"first_name || ' ' || last_name LIKE ?",
		formattedQuery,
	).Find(&users)
	err = res.Error
	return
}

func Seed() (err error) {
	users, err := GetAll()
	if err != nil {
		return
	}
	if len(users) != 0 {
		return nil
	}
	for count := 0; count < 5; count++ {
		const createdFormat = "2006-01-02 15:04:05"
		user := &userSchema.UserSchema{
			FirstName:   faker.FirstName(),
			LastName:    faker.LastName(),
			DateCreated: time.Now().Format(createdFormat),
		}
		user.Username = GenerateUsername(user.FirstName, user.LastName)

		if randomInt := mathrand.Intn(2); randomInt == 0 {
			user.DarkMode = true
		} else {
			user.DarkMode = false
		}
		_, err = Create(user)
		if err != nil {
			return
		}
	}
	return
}
