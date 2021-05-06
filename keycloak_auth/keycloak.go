package keycloak_auth

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	config Config
)

func New(cfg Config) {
	if cfg.EndpointURL.EndpointToken == "" ||
		cfg.EndpointURL.EndpointLogout == "" ||
		cfg.EndpointURL.EndpointUserInfo == "" {
		cfg.DefaultEndpointConfig()
	}
	if cfg.HashKeyFunc == nil {
		cfg.HashKeyFunc = HashSHA1Func
	}
	config = cfg
	cfg.FiberRouter.Post("/token", EndpointToken)
	cfg.FiberRouter.Get("/logout", EndpointLogout)
	cfg.FiberRouter.Get("/userinfo", EndpointUserInfo)

}

func EndpointToken(c *fiber.Ctx) error {
	urlToken := config.EndpointURL.EndpointToken
	dataForm := url.Values{}
	if grant_type := c.FormValue("grant_type", ""); grant_type == "refresh_token" {
		refresh_token := c.FormValue("refresh_token", "")
		if refresh_token == "" {
			return c.SendStatus(400)
		}
		dataForm = url.Values{
			"client_id":     []string{config.ClientID},
			"client_secret": []string{config.ClientSecret},
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refresh_token},
		}
	} else {
		sign := GetToken{}
		if err := c.BodyParser(&sign); err != nil {
			return err
		}
		if sign.Scope == "" {
			sign.Scope = "profile"
		}
		dataForm = url.Values{
			"client_id":     []string{config.ClientID},
			"client_secret": []string{config.ClientSecret},
			"grant_type":    []string{"password"},
			"username":      []string{sign.Username},
			"password":      []string{sign.Password},
			"scope":         []string{sign.Scope},
		}
	}

	resp, err := http.PostForm(urlToken, dataForm)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	respToken := ResponseToken{}
	if resp.StatusCode == 200 {
		if err := json.Unmarshal(body, &respToken); err == nil {
			config.CacheToken.Set(
				config.HashKeyFunc([]byte(respToken.AccessToken)),
				respToken, time.Duration(respToken.ExpiresIn)*time.Second,
			)
		} else {
			panic(err)
		}
	}
	c.Set("content-type", resp.Header.Get("content-type"))
	return c.Status(resp.StatusCode).Send(body)
}

func EndpointLogout(c *fiber.Ctx) error {
	headerAuth := c.Get("authorization", "")
	client := &http.Client{}
	ah := strings.Split(headerAuth, " ")
	if len(ah) < 2 {
		return c.JSON(fiber.Map{"error": "authorization header error"})
	}

	form := url.Values{}
	if a, ok := config.CacheToken.Get(config.HashKeyFunc([]byte(ah[1]))); ok {
		tk := a.(ResponseToken)
		form.Add("refresh_token", tk.RefreshToken)
	}
	form.Add("client_id", config.ClientID)
	form.Add("client_secret", config.ClientSecret)

	req, _ := http.NewRequest("POST", config.EndpointURL.EndpointLogout, strings.NewReader(form.Encode()))
	req.Header.Add("authorization", headerAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		if len(ah) == 2 {
			config.CacheToken.Delete(config.HashKeyFunc([]byte(ah[1])))
		}
	}
	c.Set("content-type", resp.Header.Get("content-type"))
	return c.Status(resp.StatusCode).Send(body)
}

func EndpointUserInfo(c *fiber.Ctx) error {
	headerAuth := c.Get("authorization", "")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.EndpointURL.EndpointUserInfo, nil)
	req.Header.Add("authorization", headerAuth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	c.Set("content-type", resp.Header.Get("content-type"))
	return c.Status(resp.StatusCode).Send(body)
}

func HashSHA1Func(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}
