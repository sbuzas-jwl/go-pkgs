package sumtype

import (
	"encoding/json"
	"errors"
	"maps"
	"slices"
)

//EX
//externally tagged: {"delete_object": {"id": "1", "soft_delete": true}}

type RegionExternal map[CountryCode]Region

func (r *RegionExternal) CountryCode() (CountryCode, error) {
	if keyLen := len(*r); keyLen == 0 {
		return "", errors.New("no country code")
	} else if keyLen == 1 {
		// 1 countryCode
		return slices.Collect(maps.Keys(*r))[0], nil
	}

	return "", errors.New("multiple-keys region invalid")

}

func (r *RegionExternal) Value() (Region, error) {
	key, err := r.CountryCode()
	if err != nil {
		return nil, err
	}

	return (*r)[key], nil
}

// TODO: This is really rough and probably fails in about a dozen cases
func (r *RegionExternal) UnmarshalJSON(data []byte) error {
	var typeHint map[CountryCode]json.RawMessage
	if err := json.Unmarshal(data, &typeHint); err != nil {
		return err
	}

	var typeKey CountryCode
	if keyLen := len(typeHint); keyLen == 1 {
		// 1 countryCode
		typeKey = slices.Collect(maps.Keys(typeHint))[0]
	}

	region, err := NewRegion(typeKey)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(typeHint[typeKey], region); err != nil {
		return err
	}

	*r = RegionExternal{
		region.CountryCode(): region,
	}
	return nil
}

func NewRegionExternal[T Region](v T) (RegionExternal, error) {
	return RegionExternal{
		v.CountryCode(): v,
	}, nil
}
