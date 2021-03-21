# Keycloakware  for `fiber v2`

Openid Connect(OIDC) with keycloak for middleware `fiber v2`

Verrify with [JSON Web Tokens](https://github.com/gofiber/jwt)
JWT returns a JSON Web Token (JWT) auth middleware.
For valid token, it sets the user in Ctx.Locals and calls next handler.
For invalid token, it returns "401 - Unauthorized" error.
For missing token, it returns "400 - Bad Request" error.

# Support 
- [x] RSA Method Only



### Install

This middleware supports Fiber v2, install accordingly.

```
go get -u github.com/attapon-th/go-pkg/fiber-middleware/keycloakware
```
**Or**
```go
import "github.com/attapon-th/go-pkg/fiber-middleware/keycloakware"
```

### Signature
```go
//  auth with Header `Authorization: Bearer <token>`
keycloakware.New(config keycloakware.Config) func(*fiber.Ctx) error
```

### Config
| Property | Type | Description | Default |
| :--- | :--- | :--- | :--- |
| BaseInstant | string | baseurl `keycloak` instant example: `http://localhost:9090/auth` | "" |
|RealmName| string | `Keycloak` realm name | "" |
| Filter | `func(*Ctx) bool` | Defines a function to skip middleware | `nil` |
| SuccessHandler | `func(*fiber.Ctx) error` |  SuccessHandler defines a function which is executed for a valid token. | `nil` |
| ErrorHandler | `func(*fiber.Ctx, error) error` | ErrorHandler defines a function which is executed for an invalid token. | `401 Invalid or expired JWT` |
| SigningKey | `interface{}` | Signing key to validate token. Used as fallback if SigningKeys has length 0. | `nil` |
| SigningKeys | `map[string]interface{}` | Map of signing keys to validate token with kid field usage. | `nil` |
| SigningMethod | `string` | Signing method, used to check token signing method. Possible values: `HS256`, `HS384`, `HS512`, `ES256`, `ES384`, `ES512`, `RS256`, `RS384`, `RS512` | `"HS256"` |
| ContextKey | `string` | Context key to store user information from the token into context. | `"user"` |
| Claims | `jwt.Claim` | Claims are extendable claims data defining token content. | `jwt.MapClaims{}` |
| TokenLookup | `string` | TokenLookup is a string in the form of `<source>:<name>` that is used | `"header:Authorization"` |
| AuthScheme | `string` |AuthScheme to be used in the Authorization header. | `"Bearer"` |

### RS512  Example
```go
package main

import (
	"fmt"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/attapon-th/go-pkg/fiber-middleware/keycloakware"
)

func main() {
	app := fiber.New()

	// Login route

	// JWT Keycloak Middleware
	app.Use("/valid", keycloakware.New(keycloakware.Config{
		BaseInstant:   "https://localhost:8000/auth",
		RealmName:     "master", 
		SigningMethod: "RS512", // fix method
	}))

	// Restricted Routes
	app.Get("/valid", restricted)

	app.Listen("localhost:3000")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
```

