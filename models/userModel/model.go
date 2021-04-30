package userModel

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/o1egl/paseto"

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
	token, err := NewAuthToken(newUser.Username)
	if err != nil {
		return
	}
	newUser.AuthToken = token
	res := db.Create(&newUser)
	user = newUser
	if err = res.Error; err != nil {
		return
	}
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

func NewAuthToken(username string) (token string, err error) {
	b, _ := hex.DecodeString("b4cbfb43df4ce210727d953e4a713307fa19bb7d9f85041438d9e11b942a37741eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	privateKey := ed25519.PrivateKey(b)

	jsonToken := paseto.JSONToken{
		Expiration: time.Now().Add(24 * time.Hour),
	}
	jsonToken.Set("username", username)
	v2 := paseto.NewV2()

	token, err = v2.Sign(privateKey, jsonToken, nil)
	if err != nil {
		return "", err
	}

	return
}

func VerifyAuthToken(token string) (user *userSchema.UserSchema, err error) {
	db, err := config.OpenDB()
	if nil != err {
		panic(err)
	}

	b, _ := hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	publicKey := ed25519.PublicKey(b)

	res := db.Where("auth_token = ?", token).First(&user)
	if err = res.Error; err != nil {
		return nil, err
	}
	v2 := paseto.NewV2()
	var newJsonToken paseto.JSONToken
	err = v2.Verify(token, publicKey, &newJsonToken, nil)
	if err != nil {
		return nil, err
	}
	return
}
