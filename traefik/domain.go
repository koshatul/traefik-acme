package traefik

import "strings"

// Domain holds a domain name with SANs.
type Domain struct {
	Main string
	SANs []string
}

// ToStrArray convert a domain into an array of strings.
func (d *Domain) ToStrArray() (domains []string) {
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

// Contains returns true if a specified domain is in the Domain object.
func (d *Domain) Contains(domain string) bool {
	if strings.EqualFold(domain, d.Main) {
		return true
	}

	if d.SANs != nil && len(d.SANs) > 0 {
		for _, san := range d.SANs {
			if strings.EqualFold(domain, san) {
				return true
			}
		}
	}

	return false
}
