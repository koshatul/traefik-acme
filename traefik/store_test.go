package traefik_test

import (
	"github.com/koshatul/traefik-acme/traefik"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

// func runAcmeDataTests(acmeDataBuf []byte) {
// 	It("should find test.example.com", func() {
// 		store, err := traefik.ReadBytes(acmeDataBuf)

// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(store).NotTo(BeNil())

// 		cert := store.GetCertificateByName("test.example.com")
// 		Expect(cert).NotTo(BeNil())
// 		Expect(cert.Domain.ToStrArray()).To(ContainElement("test.example.com"))
// 	})

// 	It("should also find another-test.example.com", func() {
// 		store, err := traefik.ReadBytes(acmeDataBuf)

// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(store).NotTo(BeNil())

// 		cert := store.GetCertificateByName("another-test.example.com")
// 		Expect(cert).NotTo(BeNil())
// 		Expect(cert.Domain.ToStrArray()).To(ContainElement("another-test.example.com"))
// 	})

// 	It("should not find test2.example.com", func() {
// 		store, err := traefik.ReadBytes(acmeDataBuf)

// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(store).NotTo(BeNil())

// 		cert := store.GetCertificateByName("test2.example.com")
// 		Expect(cert).To(BeNil())
// 	})
// }

var _ = Describe("LocalStore", func() {

	DescribeTable("should find test.example.com",
		func(acmeDataBuf *[]byte) {
			store, err := traefik.ReadBytes(*acmeDataBuf)

			Expect(err).NotTo(HaveOccurred())
			Expect(store).NotTo(BeNil())

			cert := store.GetCertificateByName("test.example.com")
			Expect(cert).NotTo(BeNil())
			Expect(cert.Domain.ToStrArray()).To(ContainElement("test.example.com"))
		},
		Entry("traefik v1", &acmeDatav1),
		Entry("traefik v2", &acmeDatav2),
	)

	DescribeTable("should also find another-test.example.com",
		func(acmeDataBuf *[]byte) {
			store, err := traefik.ReadBytes(*acmeDataBuf)

			Expect(err).NotTo(HaveOccurred())
			Expect(store).NotTo(BeNil())

			cert := store.GetCertificateByName("another-test.example.com")
			Expect(cert).NotTo(BeNil())
			Expect(cert.Domain.ToStrArray()).To(ContainElement("another-test.example.com"))
		},
		Entry("traefik v1", &acmeDatav1),
		Entry("traefik v2", &acmeDatav2),
	)

	DescribeTable("should not find test2.example.com",
		func(acmeDataBuf *[]byte) {
			store, err := traefik.ReadBytes(*acmeDataBuf)

			Expect(err).NotTo(HaveOccurred())
			Expect(store).NotTo(BeNil())

			cert := store.GetCertificateByName("test2.example.com")
			Expect(cert).To(BeNil())
		},
		Entry("traefik v1", &acmeDatav1),
		Entry("traefik v2", &acmeDatav2),
	)

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
