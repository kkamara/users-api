package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kkamara/users-api/handlers"
	"github.com/kkamara/users-api/schemas/userSchema"
)

const app_url = "http://localhost:3000/api"

var username = ""

func TestPostUser(t *testing.T) {
	app := *fiber.New()
	app.Post("/api/users", handlers.PostUser)

	user := &userSchema.UserSchema{
		FirstName: "Mary",
		LastName:  "Poppins",
		DarkMode:  false,
	}
	jsonBody, _ := json.Marshal(&user)
	req := httptest.NewRequest(
		"POST",
		fmt.Sprintf("%s/users", app_url),
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 201
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Data userSchema.UserSchema `json:"data"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
		username = res.Data.Username
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}

func TestGetUsers(t *testing.T) {
	app := *fiber.New()
	app.Get("/api/users", handlers.GetUsers)

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("%s/users", app_url),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 200
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Data []*userSchema.UserSchema `json:"data"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}

func TestSearchUsers(t *testing.T) {
	app := *fiber.New()
	app.Get("/api/users/search", handlers.SearchUsers)

	query := "Mary%20Poppins"
	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("%s/users/search?query=%s", app_url, query),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 200
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Data []*userSchema.UserSchema `json:"data"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}

func TestPutUsers(t *testing.T) {
	app := *fiber.New()
	app.Patch("/api/users/:username", handlers.PatchUser)

	user := &userSchema.UserSchema{
		FirstName: "Mary",
		LastName:  "Moppins",
		DarkMode:  true,
	}
	jsonBody, _ := json.Marshal(&user)
	req := httptest.NewRequest(
		"PATCH",
		fmt.Sprintf("%s/users/%s", app_url, username),
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 200
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Data userSchema.UserSchema `json:"data"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}

func TestToggleDarkMode(t *testing.T) {
	app := *fiber.New()
	app.Put("/api/users/:username/darkmode", handlers.PutToggleDarkMode)

	req := httptest.NewRequest(
		"PUT",
		fmt.Sprintf("%s/users/%s/darkmode", app_url, username),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 200
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Data userSchema.UserSchema `json:"data"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}

func TestDeleteUser(t *testing.T) {
	app := *fiber.New()
	app.Delete("/api/users/:username", handlers.DeleteUser)

	req := httptest.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/users/%s", app_url, username),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	expectedStatus := 200
	if resp.StatusCode == expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		var res struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Fatalf(
			`%s: StatusCode: got %v - expected %v`, t.Name(),
			resp.StatusCode,
			expectedStatus,
		)
	}
}
