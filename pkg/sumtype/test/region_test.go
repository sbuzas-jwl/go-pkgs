package test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/ca"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/mx"
	"github.com/sbuzas-jwl/go-pkgs/pkg/sumtype/regions/us"
	"github.com/stretchr/testify/assert"
)

func regionsMap() map[regions.CountryCode]regions.Region {
	var rmap = make(map[regions.CountryCode]regions.Region, 0)
	add := func(r regions.Region) {
		rmap[regions.CountryCode(r.CountryCode())] = r
	}
	usRegion := regions.New[*us.Region]()
	usRegion.SSNTail = "1234"
	usRegion.Sex = us.SexAssignedAtBirth("f")
	add(usRegion)
	caRegion := regions.New[*ca.Region]()
	caRegion.SINPrefix = "987"
	add(caRegion)
	mxRegion := regions.New[*mx.Region]()
	mxRegion.NationalID = "123456789"
	add(mxRegion)

	return rmap
}

func TestRegionInternal(t *testing.T) {
	rmap := regionsMap()
	container := func(r regions.Region) any {
		c, _ := sumtype.NewRegionInternal(r)
		return c
	}
	runRegionContainerTest[sumtype.RegionInternal](t, regions.US, container, rmap[regions.US])
	runRegionContainerTest[sumtype.RegionInternal](t, regions.CA, container, rmap[regions.CA])
	runRegionContainerTest[sumtype.RegionInternal](t, regions.MX, container, rmap[regions.MX])
}

// TODO(sbuzas): round-tripped data is not "equal" UNLESS new region adjacent is created with a pointer.
// Not sure if it's a assert.Equal limitation
func TestRegionAdjacent(t *testing.T) {
	rmap := regionsMap()
	container := func(r regions.Region) any {
		c, _ := sumtype.NewRegionAdjacent(r)
		return c
	}
	runRegionContainerTest[sumtype.RegionAdjacent](t, regions.US, container, rmap[regions.US])
	runRegionContainerTest[sumtype.RegionAdjacent](t, regions.CA, container, rmap[regions.CA])
	runRegionContainerTest[sumtype.RegionAdjacent](t, regions.MX, container, rmap[regions.MX])

}

func TestRegionExternal(t *testing.T) {
	rmap := regionsMap()
	container := func(r regions.Region) any {
		c, _ := sumtype.NewRegionExternal(r)
		return c
	}
	runRegionContainerTest[sumtype.RegionExternal](t, regions.US, container, rmap[regions.US])
	runRegionContainerTest[sumtype.RegionExternal](t, regions.US, container, rmap[regions.CA])
	runRegionContainerTest[sumtype.RegionExternal](t, regions.US, container, rmap[regions.MX])
}

func runRegionContainerTest[T any](
	t *testing.T,
	code regions.CountryCode,
	containerFn func(regions.Region) any,
	data regions.Region,
) {
	t.Run(code.String(), func(t *testing.T) {
		t.Parallel()
		container := containerFn(data)
		regionBytes, err := json.Marshal(container)
		if err != nil {
			t.Fatalf("unable to marshal [%s] region;\n value: %#v", code, container)
		}
		t.Log("Marshalled Region\n", string(regionBytes))

		var unmarshalledRegion T
		if err := json.Unmarshal(regionBytes, &unmarshalledRegion); err != nil {
			t.Fatalf("unable to unmarshal [%s] region", code)
		}

		if !cmp.Equal(container, unmarshalledRegion, cmp.AllowUnexported(sumtype.RegionInternal{})) {
			t.Fatal(cmp.Diff(container, unmarshalledRegion))
		}
	})
}

func TestRegionInternal_Embeddable(t *testing.T) {
	type EmbeddedRegionObject struct {
		Region sumtype.RegionInternal `json:"region"`
	}

	region, _ := regions.NewByCode(regions.US)
	v := region.(*us.Region)
	v.SSNTail = "1234"
	v.Sex = us.SexAssignedAtBirth("f")
	container, _ := sumtype.NewRegionInternal(region)

	obj := EmbeddedRegionObject{
		Region: container,
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
	region, _ := regions.NewByCode(regions.US)
	v := region.(*us.Region)
	v.SSNTail = "1234"
	v.Sex = us.SexAssignedAtBirth("f")
	container, _ := sumtype.NewRegionAdjacent(region)

	obj := EmbeddedRegionObject{
		Region: container,
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

	if !cmp.Equal(container, unmarshalledObj.Region) {
		t.Fatal(cmp.Diff(container, unmarshalledObj.Region))
	}
}

func TestRegionExternal_Embeddable(t *testing.T) {
	type EmbeddedRegionObject struct {
		Region sumtype.RegionExternal `json:"region"`
	}
	region, _ := regions.NewByCode(regions.US)
	v := region.(*us.Region)
	v.SSNTail = "1234"
	v.Sex = us.SexAssignedAtBirth("f")
	container, _ := sumtype.NewRegionExternal(region)

	obj := EmbeddedRegionObject{
		Region: container,
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

	if !cmp.Equal(container, unmarshalledObj.Region) {
		t.Fatal(cmp.Diff(container, unmarshalledObj.Region))
	}
}
