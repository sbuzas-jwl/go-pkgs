package us

// Typed Region specific values
type (
	SexAssignedAtBirth string
)

type Region struct {
	SSNTail string             `json:"ssn_tail"`
	Sex     SexAssignedAtBirth `json:"sex"`
}

func (r Region) CountryCode() string {
	return "us"
}
