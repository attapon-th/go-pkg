package keycloak_auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	FiberRouter    fiber.Router
	BaseInstantUrl string
	RealmName      string
	EndpointURL    EndpointConfig
	ClientID       string
	ClientSecret   string
}

type EndpointConfig struct {
	EndpointToken    string
	EndpointLogout   string
	EndpointUserInfo string
}

func (c *Config) DefaultEndpointConfig() {
	instant := strings.TrimSuffix(c.BaseInstantUrl, "/")
	realm := c.RealmName
	c.EndpointURL.EndpointToken = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", instant, realm)
	c.EndpointURL.EndpointLogout = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", instant, realm)
	c.EndpointURL.EndpointUserInfo = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", instant, realm)
}

type GetToken struct {
	Username string `json:"username" form:"username" query:"username" `
	Password string `json:"password" form:"password" query:"password"`
	Scope    string `json:"scope" form:"scope" query:"scope"`
}
