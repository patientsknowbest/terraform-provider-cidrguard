package cidr

import (
	"fmt"
	"math/big"
	"net"
)

// Range represents a range of IP addresses with its original CIDR notation and metadata.
type Range struct {
	Name        string
	CIDR        string
	Description string
	FirstIP     net.IP
	LastIP      net.IP
	Prefix      string
	Length      int
	Count       *big.Int
	ipNet       *net.IPNet
}

// NewRange creates a new Range from a CIDR string and metadata.
func NewRange(name, cidrStr, description string) (*Range, error) {
	ip, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR %q: %w", cidrStr, err)
	}

	firstIP := ip.Mask(ipNet.Mask)
	ones, bits := ipNet.Mask.Size()
	lastIP := make(net.IP, len(firstIP))
	for i := range firstIP {
		lastIP[i] = firstIP[i] | ^ipNet.Mask[i]
	}

	count := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits-ones)), nil)

	return &Range{
		Name:        name,
		CIDR:        cidrStr,
		Description: description,
		FirstIP:     firstIP,
		LastIP:      lastIP,
		Prefix:      fmt.Sprintf("%s/%d", firstIP.String(), ones),
		Length:      ones,
		Count:       count,
		ipNet:       ipNet,
	}, nil
}

// ValidateNoOverlap checks for overlapping CIDR ranges.
func ValidateNoOverlap(ranges []*Range) error {
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[i].ipNet.Contains(ranges[j].FirstIP) || ranges[i].ipNet.Contains(ranges[j].LastIP) ||
				ranges[j].ipNet.Contains(ranges[i].FirstIP) || ranges[j].ipNet.Contains(ranges[i].LastIP) {
				return fmt.Errorf("CIDR blocks for networks '%s' (%s) and '%s' (%s) overlap",
					ranges[i].Name, ranges[i].CIDR, ranges[j].Name, ranges[j].CIDR)
			}
		}
	}
	return nil
}
