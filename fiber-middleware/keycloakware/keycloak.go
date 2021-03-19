package keycloakware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

type Config struct {
	// Url Keycloak configulation well-known
	// request
	// example: `http://localhost:9090/auth`
	BaseInstant string

	// keycloak realm name
	// request
	RealmName string

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler fiber.ErrorHandler

	// Signing key to validate token. Used as fallback if SigningKeys has length 0.
	// Required. This or SigningKeys.
	SigningKey interface{}

	// Map of signing keys to validate token with kid field usage.
	// Required. This or SigningKey.
	SigningKeys map[string]interface{}

	// Signing method, used to check token signing method.
	// Optional. Default: "HS256".
	// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
	SigningMethod string

	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey string

	// Claims are extendable claims data defining token content.
	// Optional. Default value jwt.MapClaims
	Claims jwt.Claims

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default: "Bearer".
	AuthScheme string
}

func NewWithHeaderBearer(config Config) fiber.Handler {
	if config.BaseInstant == "" || config.RealmName == "" {
		panic(fmt.Errorf("Config is not empty. %v", config))
	}
	kcConfig, err := NewKeyCloakConfiguration(config.BaseInstant, config.RealmName)
	if err != nil {
		panic(err)
	}
	keyLen := len(kcConfig.Certs.Keys)
	if keyLen > 0 {
		var keys = make(map[string]interface{})
		for _, certs := range kcConfig.Certs.Keys {
			pub, err := certs.DecodePublicKey()
			if *certs.Alg == config.SigningMethod {
				keys[*certs.Kid] = pub
			}

			if err != nil {
				panic(err)
			}
		}
		return jwtware.New(jwtware.Config{
			SigningMethod:  config.SigningMethod,
			SigningKey:     config.SigningKey,
			SigningKeys:    keys,
			SuccessHandler: config.SuccessHandler,
			ErrorHandler:   config.ErrorHandler,
			TokenLookup:    config.TokenLookup,
			Filter:         config.Filter,
			ContextKey:     config.ContextKey,
			Claims:         config.Claims,
			AuthScheme:     config.AuthScheme,
		})
	} else {
		panic(fmt.Errorf("Error Create Keycloak middleware"))
	}

}
