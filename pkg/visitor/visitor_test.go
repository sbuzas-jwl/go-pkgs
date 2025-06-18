package visitor_test

import (
	"testing"

	"github.com/sbuzas-jwl/go-pkgs/pkg/visitor"
)

func Test_MemberConverison(t *testing.T) {
	elements := []visitor.Element{visitor.MemcoMember{
		ID:        "00000-000000-0000000",
		FirstName: "Bilbo",
		LastName:  "Baggins",
	}, visitor.PeoMember{
		ID:        5025600,
		FirstName: "Lotho",
		LastName:  "Sackville-Baggins",
	}, visitor.MemcoMember{
		ID:        "123456-1234-1235767",
		FirstName: "Lobelia",
		LastName:  "Sackville-Baggins",
	}}
	vis := new(visitor.MemberConversionVisitor)
	vis.ExternService = func() string { return "filthy hobbitses" }
	for _, el := range elements {
		el.Accept(vis)
	}

	t.Logf("%#v", vis.Members())

	//Test merge
	membersMap := make(map[string]visitor.Member)
	for _, m := range vis.Members() {
		membersMap[m.ID] = m
	}
	mergeVis := visitor.NewMemberMergeVisitor()
	mergeVis.Members = membersMap
	bilboOverride := visitor.Member{
		ID:        "00000-000000-0000000",
		Name:      "Bilbo Baggins",
		ExternVal: "hero of middle earth",
	}

	bilboOverride.Accept(mergeVis)

	t.Logf("%#v", mergeVis.Values())
}
