package keycloakware

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"math/big"

	"github.com/pkg/errors"
)

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// CertResponse is returned by the certs endpoint
type CertResponse struct {
	Keys []CertResponseKey `json:"keys,omitempty"`
}

type CertResponseKey struct {
	Kid     *string   `json:"kid,omitempty"`
	Kty     *string   `json:"kty,omitempty"`
	Alg     *string   `json:"alg,omitempty"`
	Use     *string   `json:"use,omitempty"`
	N       *string   `json:"n,omitempty"`
	E       *string   `json:"e,omitempty"`
	KeyOps  *[]string `json:"key_ops,omitempty"`
	X5u     *string   `json:"x5u,omitempty"`
	X5c     *[]string `json:"x5c,omitempty"`
	X5t     *string   `json:"x5t,omitempty"`
	X5tS256 *string   `json:"x5t#S256,omitempty"`
}

func (ck CertResponseKey) DecodePublicKey() (*rsa.PublicKey, error) {
	var e, n *string
	e = ck.E
	n = ck.N
	const errMessage = "could not decode public key"

	decN, err := base64.RawURLEncoding.DecodeString(*n)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	nInt := big.NewInt(0)
	nInt.SetBytes(decN)

	decE, err := base64.RawURLEncoding.DecodeString(*e)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	var eBytes []byte
	if len(decE) < 8 {
		eBytes = make([]byte, 8-len(decE), 8)
		eBytes = append(eBytes, decE...)
	} else {
		eBytes = decE
	}

	eReader := bytes.NewReader(eBytes)
	var eInt uint64
	err = binary.Read(eReader, binary.BigEndian, &eInt)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	pKey := rsa.PublicKey{N: nInt, E: int(eInt)}
	return &pKey, nil
}
