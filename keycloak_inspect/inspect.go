package keycloak_inspect

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

var (
	InspectCache *cache.Cache
	Log          *log.Logger
)

type Config struct {
	UserInfoURL    string             // require
	IntervalCache  time.Duration      // default: 1*time.Minute
	ContextKey     string             // default: user
	JsonTokenIDKey string             // default: jti
	SuccessHandler fiber.Handler      // default: success func
	ErrorHandler   fiber.ErrorHandler // default: error func
}

func New(config ...Config) fiber.Handler {
	Log = log.New(os.Stdout, "[KeycloakInspect]", log.Ldate|log.Lmicroseconds|log.Lmsgprefix)
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.UserInfoURL == "" {
		panic(fmt.Errorf("UserInfo is not empty."))
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

	InspectCache = cache.New(cfg.IntervalCache, 2*cfg.IntervalCache)

	return func(c *fiber.Ctx) error {
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
			// Log.Println(jki)
			c.Locals(cfg.JsonTokenIDKey, jki)
			if _, ok := InspectCache.Get(jki); ok {
				// Log.Println("UserInfo Cahced.")
				return cfg.SuccessHandler(c)
			}
		}
		req, err := http.NewRequest("GET", cfg.UserInfoURL, nil)
		if err == nil {
			req.Header.Add("Authorization", "Bearer "+user.Raw)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				Log.Printf("[%d]Get UserInfo from authorization server.\n", resp.StatusCode)
				if resp.StatusCode == 200 {
					b, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return c.Status(401).SendString("Token invalid: Token is not active")
					}
					// Log.Println(string(b))
					InspectCache.Set(jki, b, cfg.IntervalCache)
					return cfg.SuccessHandler(c)
				} else if resp.StatusCode == 401 {
					return c.Status(401).SendString("Token invalid: Token is not active")
				}
			}
		}
		return cfg.ErrorHandler(c, fmt.Errorf("Connection Authorization Server Error."))

	}

}
