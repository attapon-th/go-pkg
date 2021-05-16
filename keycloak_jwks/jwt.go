// Thanks to https://github.com/uandiii/go-jwk-security

package keycloak_jwks

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/attapon-th/go-pkg/keycloak_jwks/security"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	JwksCertURL    string
	KeyFunc        jwt.Keyfunc
	Claims         jwt.Claims
	SuccessHandler fiber.Handler
	ErrorHandler   fiber.ErrorHandler
}

func New(config ...Config) fiber.Handler {
	var cfg Config
	headers := "header:" + fiber.HeaderAuthorization
	authScheme := "Bearer"

	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.JwksCertURL == "" {
		panic(fmt.Errorf("JwksCertURL not empty."))
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed JWT")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
			}
		}
	}

	if cfg.KeyFunc == nil {
		jks := security.New()
		jks.JwtURL = cfg.JwksCertURL
		if err := jks.GetKeys(); err != nil {
			panic(err)
		}
		cfg.KeyFunc = jks.GetKey
	}

	parts := strings.Split(headers, ":")

	t := jwtFromHeaders(parts[1], authScheme)

	return func(c *fiber.Ctx) error {
		auth, err := t(c)

		cfg.Claims = jwt.MapClaims{}

		var token = new(jwt.Token)

		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.KeyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.KeyFunc)
		}

		if err == nil && token.Valid {
			c.Locals("user", token)
			return cfg.SuccessHandler(c)
		}

		return cfg.ErrorHandler(c, err)
	}
}

func jwtFromHeaders(header string, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)

		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}

		return "", errors.New("Missing or malformed JWT")
	}
}
