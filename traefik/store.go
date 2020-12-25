package traefik

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// LocalStore represents the parent store.
type LocalStore map[string]*LocalNamedStore

// LocalNamedStore represents the data managed by the Store.
type LocalNamedStore struct {
	// Acme           *LocalStore    `json:"acme"`
	Account        *Account       `json:"Account"`
	Certificates   []*Certificate `json:"Certificates"`
	HTTPChallenges map[string]map[string][]byte
	TLSChallenges  map[string]*Certificate
}

// ReadFile returns new LocalNamedStore from a filename.
func ReadFile(filename, certificateResolver string) (*LocalNamedStore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	defer f.Close()

	file, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return ReadBytes(file, certificateResolver)
}

// ErrCertificateResolverNotFound is returned when the specified certificate resolver is not found.
var ErrCertificateResolverNotFound = errors.New("certificate resolver not found")

// ReadBytes returns new LocalNamedStore from a byte slice.
func ReadBytes(data []byte, certificateResolver string) (*LocalNamedStore, error) {
	o := LocalStore{}
	if err := json.Unmarshal(data, &o); err != nil {
		// fallback to traefik v1 (no resolver parent key in JSON)
		v := &LocalNamedStore{}
		if err := json.Unmarshal(data, v); err == nil {
			return v, nil
		}

		return nil, fmt.Errorf("unable to parse file: %w", err)
	}

	if v, ok := o[certificateResolver]; ok {
		return v, nil
	}

	return nil, ErrCertificateResolverNotFound
}

// GetAccount returns ACME Account.
func (s *LocalNamedStore) GetAccount() *Account {
	return s.Account
}

// GetCertificates returns ACME Certificates list.
func (s *LocalNamedStore) GetCertificates() []*Certificate {
	return s.Certificates
}

// GetCertificateByName returns ACME Certificate matching supplied name.
func (s *LocalNamedStore) GetCertificateByName(name string) *Certificate {
	for _, cert := range s.GetCertificates() {
		if cert.Domain.Contains(name) {
			return cert
		}
	}

	return nil
}

// // Store is a generic interface that represents a storage.
// type Store interface {
// 	GetAccount() (*Account, error)
// 	SaveAccount(*Account) error
// 	GetCertificates() ([]*Certificate, error)
// 	SaveCertificates([]*Certificate) error

// 	GetHTTPChallengeToken(token, domain string) ([]byte, error)
// 	SetHTTPChallengeToken(token, domain string, keyAuth []byte) error
// 	RemoveHTTPChallengeToken(token, domain string) error

// 	AddTLSChallenge(domain string, cert *Certificate) error
// 	GetTLSChallenge(domain string) (*Certificate, error)
// 	RemoveTLSChallenge(domain string) error
// }
