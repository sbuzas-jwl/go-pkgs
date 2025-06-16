package regions

import (
	"fmt"

	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/ca"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/mx"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/us"
)

// Region is a simple interface implemented by all region types.
// This is meant to "loosely" tie the many regions together with one common feature, and enables their reliable encoding/decoding
type Region interface {
	CountryCode() string
}

// CountryCode is the iso-Alpha-2 country code.
type CountryCode string

func (c CountryCode) IsZero() bool {
	return c == ""
}

func (c CountryCode) String() string {
	return string(c)
}

//NOTE: Adding a CountryCode below does not magically implement it.
// A sub-package should be made for the region, and the [NewByCode] function must be exhaustively updated.

const (
	US CountryCode = "us"
	MX CountryCode = "mx"
	CA CountryCode = "ca"
)

// NewByCode creates a new instance of a region object. It will return an error if the country code is not know.
func NewByCode(code CountryCode) (Region, error) {
	var region Region
	switch code {
	case CA:
		region = New[*ca.Region]()
	case MX:
		region = New[*mx.Region]()
	case US:
		region = New[*us.Region]()
	default:
		return nil, fmt.Errorf("unknown country code [%s]", code)
	}

	return region, nil
}

type ptr[T any] interface {
	*T
}

// New creates a new instance of a region object.
func New[T ptr[U], U Region]() T {
	var u U
	return T(&u)
}
