package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

// UserCreate is a function that creates a new user.
func TestUserCreate(t *testing.T) {
	// Setup
	app := fiber.New()
	app.Post("/users", UserCreate)
	// Test
	resp, err := app.Test(httptest.NewRequest("POST", "/users", strings.NewReader(`{"username":"test","password":"test"}`)))
	// Assert

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s", string(response))
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
