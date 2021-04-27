package keycloakware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/buger/jsonparser"
)

var UrlKeycloakConfigurationFormat = `{baseInstant}/realms/{realm}/.well-known/openid-configuration`

type KeycloakConfiguration struct {
	FullURL     string
	BaseInstant string
	Realm       string
	WellKnown   []byte
	Certs       CertResponse
	ClientID    string
	ROLES       []string
}

func NewKeyCloakConfiguration(baseInstant string, realm string) (kc *KeycloakConfiguration, err error) {
	kc = new(KeycloakConfiguration)
	newFullUrl := UrlKeycloakConfigurationFormat
	newFullUrl = strings.Replace(newFullUrl, `{baseInstant}`, baseInstant, 1)
	newFullUrl = strings.Replace(newFullUrl, `{realm}`, realm, 1)
	FullUrl, err := url.ParseRequestURI(newFullUrl)
	if err != nil {
		return
	}
	kc.BaseInstant = baseInstant
	kc.Realm = realm
	kc.FullURL = FullUrl.String()

	if err = kc.getWellKnown(); err == nil {
		err = kc.getCerts()
	}
	return
}

func (kc *KeycloakConfiguration) getWellKnown() error {
	resp, err := http.Get(kc.FullURL)
	if err != nil {
		return fmt.Errorf("(getWellKnown) - %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		kc.WellKnown, err = ioutil.ReadAll(resp.Body)
	} else {
		err = fmt.Errorf(" %d: %s", resp.StatusCode, resp.Status)
	}
	if err != nil {
		return fmt.Errorf("(getWellKnown) - %s", err)
	}
	return nil
}

func (kc *KeycloakConfiguration) getCerts() error {
	url, err := jsonparser.GetString(kc.WellKnown, "jwks_uri")
	if err != nil {
		return fmt.Errorf("(getCerts:jsonparser) - %s", err.Error())
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("(getCerts:http) - %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var certs []byte
		if certs, err = ioutil.ReadAll(resp.Body); err == nil {
			err = json.Unmarshal(certs, &kc.Certs)
		}
	} else {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}
	if err != nil {
		return fmt.Errorf("(getCerts) - %s", err)
	}
	return nil
}
