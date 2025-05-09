package sumtype

import (
	"encoding/json"
)

//EX
//adjacently tagged: {"type": "delete_object", "value": {"id": "1", "soft_delete": true}}

type RegionAdjacent struct {
	Code  CountryCode `json:"code"`
	Value Region      `json:"value"`
}

func (r *RegionAdjacent) UnmarshalJSON(data []byte) error {
	var typeHint struct {
		Hint  CountryCode     `json:"code"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(data, &typeHint); err != nil {
		return err
	}

	region, err := NewRegion(typeHint.Hint)
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

func NewRegionAdjacent[T Region](v T) (RegionAdjacent, error) {
	return RegionAdjacent{
		Code:  v.CountryCode(),
		Value: v,
	}, nil
}
