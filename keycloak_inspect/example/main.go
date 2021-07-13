package example

import (
	"fmt"

	"github.com/attapon-th/go-pkg/keycloak_inspect"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(keycloak_inspect.New(keycloak_inspect.Config{
		// InspectURL: ,
	}))
	fmt.Println(app.Listen("127.0.0.1:8888"))
}
