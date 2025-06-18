package visitor

type NoopMemberVisitor struct{}

// visitEorMember implements MemberVisitor.
func (n *NoopMemberVisitor) visitEorMember(EorMember) error {
	return nil
}

// visitMember implements MemberVisitor.
func (n *NoopMemberVisitor) visitMember(Member) error {
	return nil
}

// visitMemcoMember implements MemberVisitor.
func (n *NoopMemberVisitor) visitMemcoMember(MemcoMember) error {
	return nil
}

// visitPeoMember implements MemberVisitor.
func (n *NoopMemberVisitor) visitPeoMember(PeoMember) error {
	return nil
}

var _ MemberVisitor = (*NoopMemberVisitor)(nil)
