package ca

type Region struct {
	SINPrefix string `json:"sin_prefix"`
}

func (r Region) CountryCode() string {
	return "ca"
}
