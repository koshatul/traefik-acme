package traefik

// Certificate is a struct which contains all data needed from an ACME certificate
type Certificate struct {
	Domain      Domain `json:"domain"`
	Certificate []byte `json:"certificate"`
	Key         []byte `json:"key"`
}
