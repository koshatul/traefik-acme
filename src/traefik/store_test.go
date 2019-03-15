package traefik_test

import (
	"github.com/koshatul/traefik-acme/src/traefik"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LocalStore", func() {

	It("should find test.example.com", func() {
		store, err := traefik.ReadBytes(acmeData)

		Expect(err).NotTo(HaveOccurred())
		Expect(store).NotTo(BeNil())

		cert := store.GetCertificateByName("test.example.com")
		Expect(cert).NotTo(BeNil())
		Expect(cert.Domain.ToStrArray()).To(ContainElement("test.example.com"))
	})

	It("should also find another-test.example.com", func() {
		store, err := traefik.ReadBytes(acmeData)

		Expect(err).NotTo(HaveOccurred())
		Expect(store).NotTo(BeNil())

		cert := store.GetCertificateByName("another-test.example.com")
		Expect(cert).NotTo(BeNil())
		Expect(cert.Domain.ToStrArray()).To(ContainElement("another-test.example.com"))
	})

	It("should not find test2.example.com", func() {
		store, err := traefik.ReadBytes(acmeData)

		Expect(err).NotTo(HaveOccurred())
		Expect(store).NotTo(BeNil())

		cert := store.GetCertificateByName("test2.example.com")
		Expect(cert).To(BeNil())
	})

	It("should throw an error on corrupt acme.json data", func() {
		store, err := traefik.ReadBytes([]byte("blah"))

		Expect(err).To(HaveOccurred())
		Expect(store).To(BeNil())
	})

	It("should read but not find any certs for invalid acme.json but still json", func() {
		store, err := traefik.ReadBytes([]byte("{}"))

		Expect(err).NotTo(HaveOccurred())
		Expect(store).NotTo(BeNil())

		Expect(store.GetCertificates()).To(BeEmpty())
	})

})
