package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/attapon-th/go-pkg/keycloak_inspect"
	"github.com/attapon-th/go-pkg/keycloak_jwks"
	"github.com/gofiber/fiber/v2"
)

var (
	KcJwkConfig     keycloak_jwks.Config
	KcInspectConfig keycloak_inspect.Config
)

func main() {
	app := fiber.New()

	KcInspectConfig.UserInfoURL = os.Getenv("USER_INFO_URL")
	KcJwkConfig.JwksCertURL = os.Getenv("JWK_CERT_URL")

	app.Use(
		keycloak_jwks.New(KcJwkConfig),

		keycloak_inspect.New(KcInspectConfig),
	)
	app.Get("/", func(c *fiber.Ctx) error {
		jti, ok := c.Locals("jti").(string)
		// log.Println(jti)
		if !ok {
			return c.SendStatus(500)
		}
		j, _ := keycloak_inspect.InspectCache.Get(jti)
		info := make(map[string]interface{})
		json.Unmarshal(j.([]byte), &info)
		f := fiber.Map{
			"jti":      jti,
			"userinfo": info,
		}
		return c.JSON(f)
	})
	fmt.Println(app.Listen("127.0.0.1:8888"))
}
