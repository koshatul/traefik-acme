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

	acmeTemp := traefik.LocalStore{
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
		Acme: &acmeTemp,
	}
	err = json.NewEncoder(buf).Encode(acmeTempv2)
	Expect(err).NotTo(HaveOccurred())

	acmeDatav2 = buf.Bytes()
	Expect(acmeDatav2).NotTo(BeEmpty())
})
