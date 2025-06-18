package visitor

import (
	"maps"
	"slices"
)

type memberMergeStrategy func(*Member, Member)

func OverwriteStrategy(base *Member, new Member) {
	base.ExternVal = new.ExternVal
	base.Name = new.Name
}

type MemberMergeVisitor struct {
	NoopMemberVisitor
	Members map[string]Member
	MergeFn memberMergeStrategy
}

func NewMemberMergeVisitor() *MemberMergeVisitor {
	return &MemberMergeVisitor{
		Members: make(map[string]Member, 0),
		MergeFn: OverwriteStrategy,
	}
}

func (mv *MemberMergeVisitor) Values() []Member {
	return slices.Collect(maps.Values(mv.Members))
}

// visitMember implements MemberVisitor.
func (mv *MemberMergeVisitor) visitMember(o Member) error {
	base, ok := mv.Members[o.ID]
	if !ok {
		return nil
	}

	mv.MergeFn(&base, o)

	mv.Members[base.ID] = base
	return nil
}

var _ MemberVisitor = (*MemberMergeVisitor)(nil)
