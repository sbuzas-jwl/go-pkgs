package sumtype

//sumtype:decl
type IRegion interface {
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
	MX             = "mx"
	CA             = "ca"
)

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
