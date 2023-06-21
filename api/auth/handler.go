package auth

import (
	"apisrv/database"
	"apisrv/pkg"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	database.DefaultModel
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserCreate(c *fiber.Ctx) error {

	user := User{}
	if err := c.BodyParser(&user); err != nil {
		return pkg.BadRequest(fmt.Sprintf("failed to parse request body: %v", err))
	}

	if user.Username == "" || user.Password == "" {
		return pkg.BadRequest("username and password are required")
	}

	if err := user.Create(); err != nil {
		return pkg.Unexpected(fmt.Sprintf("failed to create user: %v", err))
	}

	resp := pkg.NewRes().Omit("Password").Transform(user)

	return c.JSON(resp)
}

func FindUser(username string) (User, error) {
	var user User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (u *User) Create() error {
	password, err := pkg.GeneratePasswordHash(u.Password)

	if err != nil {
		return err
	}

	u.Password = password
	return database.DB.Create(u).Error
}
