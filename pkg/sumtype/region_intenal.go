package sumtype

import (
	"encoding/json"

	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions"
)

// RegionInternal is an internally tagged Region sumtype wrapper.
type RegionInternal struct {
	payload json.RawMessage
}

func (a *RegionInternal) CountryCode() (regions.CountryCode, error) {
	var region struct {
		Code regions.CountryCode `json:"code"`
	}

	err := json.Unmarshal(a.payload, &region)

	return region.Code, err
}

func (r *RegionInternal) Value() (regions.Region, error) {
	code, err := r.CountryCode()
	if err != nil {
		return nil, err
	}

	v, err := regions.NewByCode(code)
	if err != nil {
		return nil, err
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

func NewRegionInternal[T regions.Region](v T) (RegionInternal, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return RegionInternal{}, err
	}

	data, err = addCountryCode(data, regions.CountryCode(v.CountryCode()))
	if err != nil {
		return RegionInternal{}, err
	}

	return RegionInternal{
		payload: data,
	}, nil
}

func addCountryCode(data json.RawMessage, code regions.CountryCode) (json.RawMessage, error) {
	var discriminatedType map[string]any
	if err := json.Unmarshal(data, &discriminatedType); err != nil {
		return nil, err
	}

	discriminatedType["code"] = code
	return json.Marshal(&discriminatedType)
}
