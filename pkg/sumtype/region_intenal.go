package sumtype

import (
	"encoding/json"
	"fmt"
)

// RegionInternal is an internally tagged Region sumtype wrapper.
type RegionInternal struct {
	payload json.RawMessage
}

func (a *RegionInternal) CountryCode() (CountryCode, error) {
	var region struct {
		Code CountryCode `json:"code"`
	}

	err := json.Unmarshal(a.payload, &region)

	return region.Code, err
}

func (r *RegionInternal) Value() (any, error) {
	code, err := r.CountryCode()
	if err != nil {
		return nil, err
	}

	var v any
	switch code {
	case US:
		v = new(USRegion)
	case MX:
		v = new(MXRegion)
	case CA:
		v = new(CARegion)
	default:
		return nil, fmt.Errorf("invalid region type: %s", code)
	}

	err = json.Unmarshal(r.payload, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r RegionInternal) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(r.payload)

	return data, err
}

func (r *RegionInternal) UnmarshalJSON(data []byte) error {
	err := r.payload.UnmarshalJSON(data)
	return err
}

func NewRegionInternal[T Region](v T) (RegionInternal, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return RegionInternal{}, err
	}

	data, err = addCountryCode(data, v.CountryCode())
	if err != nil {
		return RegionInternal{}, err
	}

	return RegionInternal{
		payload: data,
	}, nil
}

func addCountryCode(data json.RawMessage, code CountryCode) (json.RawMessage, error) {
	var discriminatedType map[string]any
	if err := json.Unmarshal(data, &discriminatedType); err != nil {
		return nil, err
	}

	discriminatedType["code"] = code
	return json.Marshal(&discriminatedType)
}
