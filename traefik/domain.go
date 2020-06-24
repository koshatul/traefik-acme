package traefik

// Domain holds a domain name with SANs.
type Domain struct {
	Main string
	SANs []string
}

// ToStrArray convert a domain into an array of strings.
func (d *Domain) ToStrArray() []string {
	var domains []string
	if len(d.Main) > 0 {
		domains = []string{d.Main}
	}
	return append(domains, d.SANs...)
}

// Set sets a domains from an array of strings.
func (d *Domain) Set(domains []string) {
	if len(domains) > 0 {
		d.Main = domains[0]
		d.SANs = domains[1:]
	}
}
