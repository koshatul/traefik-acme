package traefik_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/koshatul/traefik-acme/traefik"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// nolint: gochecknoglobals // acmeDatav1 is a test variable
var acmeDatav1 []byte

// nolint: gochecknoglobals // acmeDatav2 is a test variable
var acmeDatav2 []byte

// nolint: gochecknoglobals,lll // acmeDatav3 is a test variable
var acmeDatav3 []byte = []byte(`{"acme":{"Account":{"Email":"koshatul@noreply.users.github.com","Registration":{"body":{"status":"valid","contact":["mailto:koshatul@noreply.users.github.com"]},"uri":"https://acme-v02.api.letsencrypt.org/acme/acct/123456789"},"PrivateKey":"c2VjcmV0LXByaXZhdGUta2V5LWZvci0xMjM0NTY3ODkK","KeyType":"4096"},"Certificates":[{"domain":{"main":"example.com","sans":["*.example.com"]},"certificate":"Y2VydGlmaWNhdGUtZm9yLWV4YW1wbGUuY29tCg==","key":"a2V5LWZvci1leGFtcGxlLmNvbQo=","Store":"default"}]}}`)

// nolint: gochecknoglobals,lll // acmeDatav4 is a test variable
var acmeDatav4 []byte = []byte(`{"acme":{"Account":{"Email":"koshatul@noreply.users.github.com","Registration":{"body":{"status":"valid","contact":["mailto:koshatul@noreply.users.github.com"]},"uri":"https://acme-v02.api.letsencrypt.org/acme/acct/123456789"},"PrivateKey":"c2VjcmV0LXByaXZhdGUta2V5LWZvci0xMjM0NTY3ODkK","KeyType":"4096"},"Certificates":[{"domain":{"main":"*.example.com"},"certificate":"Y2VydGlmaWNhdGUtZm9yLWV4YW1wbGUuY29tCg==","key":"a2V5LWZvci1leGFtcGxlLmNvbQo=","Store":"default"}]}}`)

// nolint: gochecknoglobals,lll // acmeDatav5 is a test variable
var acmeDatav5 []byte = []byte(`{"acme-different":{"Account":{"Email":"koshatul@noreply.users.github.com","Registration":{"body":{"status":"valid","contact":["mailto:koshatul@noreply.users.github.com"]},"uri":"https://acme-v02.api.letsencrypt.org/acme/acct/123456789"},"PrivateKey":"c2VjcmV0LXByaXZhdGUta2V5LWZvci0xMjM0NTY3ODkK","KeyType":"4096"},"Certificates":[{"domain":{"main":"example.com","sans":["*.example.com"]},"certificate":"Y2VydGlmaWNhdGUtZm9yLWV4YW1wbGUuY29tCg==","key":"a2V5LWZvci1leGFtcGxlLmNvbQo=","Store":"default"}]}}`)

var _ = BeforeSuite(func() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	Expect(err).NotTo(HaveOccurred())

	notBefore := time.Now().Add(-time.Hour)
	notAfter := time.Now().Add(time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	Expect(err).NotTo(HaveOccurred())

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		DNSNames:  []string{"test.example.com", "another-test.example.com"},

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	Expect(err).NotTo(HaveOccurred())

	certBuf := &bytes.Buffer{}
	keyBuf := &bytes.Buffer{}

	err = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	Expect(err).NotTo(HaveOccurred())
	err = pem.Encode(keyBuf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	Expect(err).NotTo(HaveOccurred())

	acmeTemp := traefik.LocalNamedStore{
		Certificates: []*traefik.Certificate{
			{
				Key:         keyBuf.Bytes(),
				Certificate: certBuf.Bytes(),
				Domain: traefik.Domain{
					Main: "test.example.com",
					SANs: []string{"another-test.example.com"},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(acmeTemp)
	Expect(err).NotTo(HaveOccurred())

	acmeDatav1 = buf.Bytes()
	Expect(acmeDatav1).NotTo(BeEmpty())

	buf.Reset()

	acmeTempv2 := traefik.LocalStore{
		"acme": &acmeTemp,
	}
	err = json.NewEncoder(buf).Encode(acmeTempv2)
	Expect(err).NotTo(HaveOccurred())

	acmeDatav2 = buf.Bytes()
	Expect(acmeDatav2).NotTo(BeEmpty())
})
