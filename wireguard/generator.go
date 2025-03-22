package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/curve25519"
)

func GenerateKeys() (string, string, error) {
	var priv [32]byte
	_, err := rand.Read(priv[:])
	if err != nil {
		return "", "", fmt.Errorf("error generating private key: %v", err)
	}

	pub, err := curve25519.X25519(priv[:], curve25519.Basepoint)
	if err != nil {
		return "", "", errors.New("failed to generate public key")
	}

	privateKey := base64.StdEncoding.EncodeToString(priv[:])
	publicKey := base64.StdEncoding.EncodeToString(pub)

	return privateKey, publicKey, nil
}
