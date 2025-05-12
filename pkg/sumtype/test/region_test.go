package test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype"
	"github.com/stretchr/testify/assert"
)

func TestRegionInternal(t *testing.T) {
	us, _ := sumtype.NewRegionInternal(sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
	})
	runRegionTest(t, sumtype.US, us)
	ca, _ := sumtype.NewRegionInternal(sumtype.CARegion{
		SINPrefix: "987",
	})
	runRegionTest(t, sumtype.CA, ca)
	mx, _ := sumtype.NewRegionInternal(sumtype.MXRegion{
		NationalID: "123456789",
	})
	runRegionTest(t, sumtype.MX, mx)

}

// TODO(sbuzas): round-tripped data is not "equal" UNLESS new region adjacent is created with a pointer.
// Not sure if it's a assert.Equal limitation
func TestRegionAdjacent(t *testing.T) {
	us, _ := sumtype.NewRegionAdjacent(&sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
	})
	runRegionTest(t, sumtype.US, us)
	ca, _ := sumtype.NewRegionAdjacent(&sumtype.CARegion{
		SINPrefix: "987",
	})
	runRegionTest(t, sumtype.CA, ca)
	mx, _ := sumtype.NewRegionAdjacent(&sumtype.MXRegion{
		NationalID: "123456789",
	})
	runRegionTest(t, sumtype.MX, mx)

}

func TestRegionExternal(t *testing.T) {
	us, _ := sumtype.NewRegionExternal(&sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
	})
	runRegionTest(t, sumtype.US, us)
	ca, _ := sumtype.NewRegionExternal(&sumtype.CARegion{
		SINPrefix: "987",
	})
	runRegionTest(t, sumtype.CA, ca)
	mx, _ := sumtype.NewRegionExternal(&sumtype.MXRegion{
		NationalID: "123456789",
	})
	runRegionTest(t, sumtype.MX, mx)

}

func runRegionTest[T any](t *testing.T, code sumtype.CountryCode, region T) {
	t.Run(code.String(), func(t *testing.T) {
		t.Parallel()
		regionBytes, err := json.Marshal(region)
		if err != nil {
			t.Fatalf("unable to marshal [%s] region;\n value: %#v", code, region)
		}
		t.Log("Marshalled Region\n", string(regionBytes))

		var unmarshalledRegion T
		if err := json.Unmarshal(regionBytes, &unmarshalledRegion); err != nil {
			t.Fatalf("unable to unmarshal [%s] region", code)
		}

		if !cmp.Equal(region, unmarshalledRegion, cmp.AllowUnexported(sumtype.RegionInternal{})) {
			t.Fatal(cmp.Diff(region, unmarshalledRegion))
		}
	})
}

func TestRegionInternal_Embeddable(t *testing.T) {
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

	assert.EqualValues(t, obj, unmarshalledObj)
}

func TestRegionAdjacent_Embeddable(t *testing.T) {
	type EmbeddedRegionObject struct {
		Region sumtype.RegionAdjacent `json:"region"`
	}
	us, _ := sumtype.NewRegionAdjacent(&sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
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

	if !cmp.Equal(us, unmarshalledObj.Region) {
		t.Fatal(cmp.Diff(us, unmarshalledObj.Region))
	}
}

func TestRegionExternal_Embeddable(t *testing.T) {
	type EmbeddedRegionObject struct {
		Region sumtype.RegionExternal `json:"region"`
	}
	us, _ := sumtype.NewRegionExternal(&sumtype.USRegion{
		SSNTail: "1234",
		Sex:     "f",
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

	if !cmp.Equal(us, unmarshalledObj.Region) {
		t.Fatal(cmp.Diff(us, unmarshalledObj.Region))
	}
}
