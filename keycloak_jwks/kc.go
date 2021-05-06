package keycloak_jwks

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/juandiii/go-jwk-security/v2/jwt"
	"github.com/juandiii/go-jwk-security/v2/security"
)

type Config struct {
	BaseInstantUrl  string
	RealmName       string
	EndpointCertUrl string
}

func (c *Config) FixEndpointURL() {
	c.EndpointCertUrl = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/cert", c.BaseInstantUrl, c.RealmName)
}

func New(config Config) fiber.Handler {
	if config.EndpointCertUrl == "" {
		config.FixEndpointURL()
	}

	jwtConfig := &security.JwtKeys{JwtURL: config.EndpointCertUrl}
	return jwt.JwtMiddleware(jwt.Config{
		KeyFunc: jwtConfig.GetKey(),
	})
}
