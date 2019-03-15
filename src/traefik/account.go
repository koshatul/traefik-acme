package traefik

import (
	"context"
	"crypto"
	"crypto/x509"

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"
)

// Account is used to store lets encrypt registration info
type Account struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   []byte
	KeyType      certcrypto.KeyType
}

// GetEmail returns email
func (a *Account) GetEmail() string {
	return a.Email
}

// GetRegistration returns lets encrypt registration resource
func (a *Account) GetRegistration() *registration.Resource {
	return a.Registration
}

// GetPrivateKey returns private key
func (a *Account) GetPrivateKey() crypto.PrivateKey {
	privateKey, err := x509.ParsePKCS1PrivateKey(a.PrivateKey)
	if err != nil {
		return nil
	}

	return privateKey
}

// GetKeyType used to determine which algo to used
func GetKeyType(ctx context.Context, value string) certcrypto.KeyType {
	switch value {
	case "EC256":
		return certcrypto.EC256
	case "EC384":
		return certcrypto.EC384
	case "RSA2048":
		return certcrypto.RSA2048
	case "RSA4096":
		return certcrypto.RSA4096
	case "RSA8192":
		return certcrypto.RSA8192
	case "":
		return certcrypto.RSA4096
	default:
		return certcrypto.RSA4096
	}
}
