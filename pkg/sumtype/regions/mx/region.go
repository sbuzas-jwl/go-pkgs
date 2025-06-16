package mx

type Region struct {
	NationalID string `json:"national_id"`
}

func (r Region) CountryCode() string {
	return "mx"
}
