package keycloak_inspect

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	InspectURL     string
	ClientID       string
	ClientSecret   string
	IntervalCache  time.Duration
	ContextKey     string
	JsonTokenIDKey string
	SuccessHandler fiber.Handler
	ErrorHandler   fiber.ErrorHandler
	inspectCache   *cache.Cache
}

func New(config ...Config) fiber.Handler {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.InspectURL == "" {
		panic(fmt.Errorf("InspectURL is not empty."))
	}

	if cfg.ContextKey == "" {
		// get from `fiber.Local`
		// default set by `https://github.com/gofiber/jwt/v2`
		cfg.ContextKey = "user"
	}

	if cfg.IntervalCache == 0 {
		cfg.IntervalCache = time.Minute
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
	}

	if cfg.JsonTokenIDKey == "" {
		cfg.JsonTokenIDKey = "jti"
	}

	cfg.inspectCache = cache.New(cfg.IntervalCache, 2*cfg.IntervalCache)

	return func(c *fiber.Ctx) error {
		var data url.Values
		var jki string
		user := c.Locals(cfg.ContextKey).(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if err := claims.Valid(); err != nil {
			return cfg.ErrorHandler(c, err)
		}
		if jj, ok := claims[cfg.JsonTokenIDKey]; !ok || jj.(string) == "" {
			return cfg.ErrorHandler(c, fmt.Errorf("JsonTokenIDKey is not found."))
		} else {
			jki = jj.(string)
			if _, ok := cfg.inspectCache.Get(jki); ok {
				return cfg.SuccessHandler(c)
			}
		}
		data.Add("token_type_hint", "requesting_party_token")
		data.Add("token", user.Raw)
		req, _ := http.NewRequest("POST", cfg.InspectURL, strings.NewReader(data.Encode()))
		req.Header.Add("Authorization", "Basic "+basicAuth(cfg.ClientID, cfg.ClientSecret))
		resp, err := http.PostForm(cfg.InspectURL, data)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					var rd = JTIActive{}
					if json.Unmarshal(body, &rd) == nil {
						if rd.Active {
							cfg.inspectCache.Set(jki, true, cfg.IntervalCache)
							return cfg.SuccessHandler(c)
						}
					}
				}
			}
		}
		return cfg.ErrorHandler(c, fmt.Errorf("Connection Authorization Server Error."))

	}

}

type JTIActive struct {
	Active bool `json:"active"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
