package security

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

type JwtKeys struct {
	Ctx       context.Context
	JwtURL    string
	cachedSet jwk.Set
}

func New() *JwtKeys {
	jk := new(JwtKeys)
	jk.Ctx = context.Background()
	return jk
}

func (j *JwtKeys) GetKeys() error {
	if j.cachedSet != nil {
		return nil
	}

	fmt.Println("Connecting :: Keycloak")

	set, err := jwk.Fetch(j.Ctx, j.JwtURL)

	if err != nil {
		return errors.New("Couldn't connect Keycloak, try again")
	}
	fmt.Println("Connected successfully :: Keycloak")
	j.cachedSet = set

	return nil
}

func (j *JwtKeys) GetKey(token *jwt.Token) (interface{}, error) {

	keyID, ok := token.Header["kid"].(string)

	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}
	if key, ok := j.cachedSet.LookupKeyID(keyID); ok {
		var raw interface{}
		return raw, key.Raw(&raw)
	}

	return nil, fmt.Errorf("unable to find key %q", keyID)
}
