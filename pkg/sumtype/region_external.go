package sumtype

import (
	"encoding/json"
	"maps"
	"slices"
)

//EX
//externally tagged: {"delete_object": {"id": "1", "soft_delete": true}}

type RegionExternal map[CountryCode]Region

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
