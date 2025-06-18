package visitor

import (
	"fmt"
	"strconv"
)

type MemberConversionVisitor struct {
	members       []Member
	ExternService func() string
}

var _ MemberVisitor = (*MemberConversionVisitor)(nil)

// visitMemcoMember implements Visitor.
func (m *MemberConversionVisitor) visitMember(in Member) error {
	m.members = append(m.members, in)
	return nil
}

// visitMemcoMember implements Visitor.
func (m *MemberConversionVisitor) visitMemcoMember(in MemcoMember) error {
	m.members = append(m.members, Member{
		ID:        in.ID,
		Name:      fmt.Sprintf("%s %s", in.FirstName, in.LastName),
		ExternVal: m.externService(),
	})
	return nil
}

// visitPeoMember implements Visitor.
func (m *MemberConversionVisitor) visitPeoMember(in PeoMember) error {
	m.members = append(m.members, Member{
		ID:        strconv.FormatInt(int64(in.ID), 10),
		Name:      fmt.Sprintf("%s %s", in.FirstName, in.LastName),
		ExternVal: m.externService(),
	})
	return nil
}

// visitEorMember implements Visitor.
func (m *MemberConversionVisitor) visitEorMember(in EorMember) error {
	m.members = append(m.members, Member{
		ID:        strconv.FormatInt(int64(in.ID), 10),
		Name:      fmt.Sprintf("%s %s", in.FirstName, in.LastName),
		ExternVal: m.externService(),
	})
	return nil
}

func (m *MemberConversionVisitor) externService() string {
	if m.ExternService == nil {
		return ""
	}
	return m.ExternService()
}

func (m *MemberConversionVisitor) Members() []Member {
	return m.members
}
