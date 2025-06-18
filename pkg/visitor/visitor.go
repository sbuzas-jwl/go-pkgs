package visitor

type MemberVisitor interface {
	visitMemcoMember(MemcoMember) error
	visitPeoMember(PeoMember) error
	visitEorMember(EorMember) error
	visitMember(Member) error
}

type Element interface {
	Accept(v MemberVisitor) error
}
