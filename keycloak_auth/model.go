package keycloak_auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Cacher interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, Expire time.Duration)
	Delete(key string)
}

type Config struct {
	FiberRouter    fiber.Router
	BaseInstantUrl string
	RealmName      string
	EndpointURL    EndpointConfig
	ClientID       string
	ClientSecret   string
	CacheToken     Cacher
	HashKeyFunc    func(b []byte) string
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

type ResponseToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope"`
}
