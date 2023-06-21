package middleware

import (
	"apisrv/api/auth"
	"apisrv/pkg"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

var Secret string

func useAuth(app *fiber.App) {

	Secret = os.Getenv("JWT_SECRET")
	if Secret == "" {
		panic("JWT_SECRET is not set")
	}

	// Login route
	app.Post("/api/login", login)
	// Unauthenticated route
	app.Get("/", accessible)

	app.Get("/api/check", chechStatus)

	// Routes that don't need authentication
	filter := []string{
		"/api/login",
		"/",
		"/api/restricted",
		"/api/weather",
		"/api/oli",
		"/api/ip",
		"/api/phone",
		"/api/kuaidi",
		"/api/speak",
		"/api/history",
		"/api/stat",
		"/api/config",
		"/api/dream/hot-search",
		"/api/dream/category",
		"/api/dream/tags",
		"/api/dream/search",
		"/api/cookbook/recommend",
		"/api/cookbook/search",
		"/api/cookbook/detail",
		"/api/cookbook/hotSearch",
		"/api/cookbook/category",
		"/api/cookbook/cookbooks",
		"/api/cookbook/crawl",
	}

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(Secret),
		Filter: func(c *fiber.Ctx) bool {
			for _, v := range filter {
				if v == c.Path() {
					return true
				}
			}
			return false
		},
	}))

	// Restricted Routes
	app.Get("/api/restricted", restricted)

}

type authReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(c *fiber.Ctx) error {

	// Parse JSON input
	var data authReq
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	name := data.Username
	pass := data.Password

	user, err := auth.FindUser(name)
	if err != nil || pkg.ComparePasswordHash(user.Password, pass) != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  name,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(Secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"name": name, "id": user.ID, "admin": true, "token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok || user == nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
		})
	}
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func chechStatus(c *fiber.Ctx) error {

	// read token from header
	header := c.GetReqHeaders()

	token := strings.Split(header["Authorization"], " ")[1]

	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
		})
	}
	// check token
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
		})
	}
	// get claims
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
		})
	}
	// check expire
	exp := claims["exp"].(float64)
	if exp < float64(time.Now().Unix()) {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Unauthorized",
		})
	}

	return c.SendString(claims["name"].(string) + " is login")
}
