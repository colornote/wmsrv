package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/utils"
)

func Setup(app *fiber.App) {

	if os.Getenv("ENABLE_CSRF") != "" {
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-Csrf-Token",
			CookieName:     "csrf_",
			CookieSameSite: "Strict",
			Expiration:     1 * time.Hour,
			KeyGenerator:   utils.UUID,
			Extractor: func(c *fiber.Ctx) (string, error) {
				fmt.Printf("Error csrf token: %v\n", c.GetReqHeaders())
				return c.FormValue("csrf_token"), nil
			},
		}))

	}

	if os.Getenv("ENABLE_LIMITER") != "" {

		app.Use(limiter.New(limiter.Config{
			Max:               120,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.FixedWindow{},
		}))
	}

	if os.Getenv("ENABLE_LOGGER") != "" {

		app.Use(logger.New(logger.Config{
			Format:     "[${time}] ${ips} ${status} - ${latency}  ${method} ${path}?${queryParams}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "China/ShangHai",
		}))
	}

	if os.Getenv("ENABLE_CACHE") != "" {

		app.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.Query("refresh") == "true"
			},
			Expiration:   30 * time.Minute,
			CacheControl: true,
		}))
	}

	// Or extend your config for customization
	app.Use("/public", filesystem.New(filesystem.Config{
		Root:         http.Dir("./assets"),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(etag.New())

	if os.Getenv("ENABLE_MONITOR") != "" {
		app.Use(pprof.New())
		app.Get("api/dashboard", monitor.New())
	}

	useAuth(app)

}
