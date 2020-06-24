package traefik

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

// LocalStore represents the data managed by the Store
type LocalStore struct {
	Acme           *LocalStore    `json:"acme"`
	Account        *Account       `json:"Account"`
	Certificates   []*Certificate `json:"Certificates"`
	HTTPChallenges map[string]map[string][]byte
	TLSChallenges  map[string]*Certificate
}

// ReadFile returns new LocalStore from a filename
func ReadFile(filename string) (*LocalStore, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ReadBytes(file)
}

// ReadBytes returns new LocalStore from a byte slice
func ReadBytes(data []byte) (*LocalStore, error) {
	o := &LocalStore{}
	if err := json.Unmarshal(data, o); err != nil {
		return nil, err
	}

	if o.Acme != nil {
		o = o.Acme
	}

	return o, nil
}

// GetAccount returns ACME Account
func (s *LocalStore) GetAccount() *Account {
	return s.Account
}

// GetCertificates returns ACME Certificates list
func (s *LocalStore) GetCertificates() []*Certificate {
	return s.Certificates
}

// GetCertificateByName returns ACME Certificate matching supplied name
func (s *LocalStore) GetCertificateByName(name string) *Certificate {
	for _, cert := range s.GetCertificates() {
		certDomains := cert.Domain.ToStrArray()
		sort.Strings(certDomains)

		i := sort.SearchStrings(certDomains, name)
		if i < len(certDomains) && certDomains[i] == name {
			return cert
		}
	}

	return nil
}

// // Store is a generic interface that represents a storage
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
