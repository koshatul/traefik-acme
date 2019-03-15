package traefik

// Certificate is a struct which contains all data needed from an ACME certificate
type Certificate struct {
	Domain      Domain
	Certificate []byte
	Key         []byte
}
