package sumtype

import (
	"encoding/json"

	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions"
)

//EX
//adjacently tagged: {"type": "delete_object", "value": {"id": "1", "soft_delete": true}}

type RegionAdjacent struct {
	Code  regions.CountryCode `json:"code"`
	Value regions.Region      `json:"value"`
}

func (r *RegionAdjacent) UnmarshalJSON(data []byte) error {
	var typeHint struct {
		Hint  regions.CountryCode `json:"code"`
		Value json.RawMessage     `json:"value"`
	}
	if err := json.Unmarshal(data, &typeHint); err != nil {
		return err
	}

	region, err := regions.NewByCode(typeHint.Hint)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(typeHint.Value, region); err != nil {
		return err
	}

	*r = RegionAdjacent{
		Code:  typeHint.Hint,
		Value: region,
	}
	return nil
}

func NewRegionAdjacent[T regions.Region](v T) (RegionAdjacent, error) {
	return RegionAdjacent{
		Code:  regions.CountryCode(v.CountryCode()),
		Value: v,
	}, nil
}
