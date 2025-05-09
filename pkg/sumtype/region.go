package sumtype

import "fmt"

type Region interface {
	CountryCode() CountryCode
	sealed()
}

type CountryCode string

func (c CountryCode) IsZero() bool {
	return c == ""
}

func (c CountryCode) String() string {
	return string(c)
}

const (
	US CountryCode = "us"
	MX CountryCode = "mx"
	CA CountryCode = "ca"
)

func NewRegion(code CountryCode) (Region, error) {
	var region Region
	switch code {
	case CA:
		region = new(CARegion)
	case MX:
		region = new(MXRegion)
	case US:
		region = new(USRegion)
	default:
		return nil, fmt.Errorf("unknown country code [%s]", code)
	}

	return region, nil
}

type USRegion struct {
	SSNTail string `json:"ssn_tail"`
	Sex     string `json:"sex"`
}

func (r USRegion) CountryCode() CountryCode {
	return US
}

type MXRegion struct {
	NationalID string `json:"national_id"`
}

func (r MXRegion) CountryCode() CountryCode {
	return MX
}

type CARegion struct {
	SINPrefix string `json:"sin_prefix"`
}

func (r CARegion) CountryCode() CountryCode {
	return CA
}

func (r USRegion) sealed() {}
func (r MXRegion) sealed() {}
func (r CARegion) sealed() {}
