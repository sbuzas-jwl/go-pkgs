package test

import (
	"encoding/json"
	"testing"

	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype"
	"github.com/stretchr/testify/assert"
)

func TestRegion(t *testing.T) {
	us, _ := sumtype.NewRegionInternal(sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
	})
	runRegionTest(t, us)
	ca, _ := sumtype.NewRegionInternal(sumtype.CARegion{
		SINPrefix: "987",
	})
	runRegionTest(t, ca)
	mx, _ := sumtype.NewRegionInternal(sumtype.MXRegion{
		NationalID: "123456789",
	})
	runRegionTest(t, mx)

}

func runRegionTest(t *testing.T, region sumtype.RegionInternal) {
	code, err := region.CountryCode()
	if code.IsZero() || err != nil {
		t.Error("region country code is zero")
		return
	}
	t.Run(code.String(), func(t *testing.T) {
		t.Parallel()
		regionBytes, err := json.Marshal(region)
		if err != nil {
			t.Fatalf("unable to marshal [%s] region;\n value: %#v", code, region)
		}
		t.Log("Marshalled Region\n", string(regionBytes))

		var unmarshalledRegion sumtype.RegionInternal
		if err := json.Unmarshal(regionBytes, &unmarshalledRegion); err != nil {
			t.Fatalf("unable to unmarshal [%s] region", code)
		}

		assert.Equal(t, region, unmarshalledRegion)
	})

}

func TestRegion_Embeddable(t *testing.T) {
	type EmbeddedRegionObject struct {
		Region sumtype.RegionInternal `json:"region"`
	}
	us, _ := sumtype.NewRegionInternal(sumtype.USRegion{
		SSNTail: "1234",
	})

	obj := EmbeddedRegionObject{
		Region: us,
	}

	objBytes, err := json.Marshal(obj)
	if err != nil {
		t.Fatal("unable to marshal", err)
	}

	t.Log("Marshalled Embeddable Region\n", string(objBytes))

	var unmarshalledObj EmbeddedRegionObject
	if err := json.Unmarshal(objBytes, &unmarshalledObj); err != nil {
		t.Fatal("unable to unmarshal", err)
	}

	assert.Equal(t, obj, unmarshalledObj)
}
