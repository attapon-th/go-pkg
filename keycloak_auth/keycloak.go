package keycloak_auth

import (
	"io/ioutil"
	"net/http"
	"net/url"

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
	config = cfg
	cfg.FiberRouter.Post("/token", EndpointToken)
	cfg.FiberRouter.Post("/logout", EndpointLogout)
	cfg.FiberRouter.Post("/userinfo", EndpointUserInfo)

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
	c.Set("content-type", resp.Header.Get("content-type"))
	return c.Send(body)
}

func EndpointLogout(c *fiber.Ctx) error {
	headerAuth := c.Get("authorization", "")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.EndpointURL.EndpointLogout, nil)
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
	return c.Send(body)
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
	return c.Send(body)
}
