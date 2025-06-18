package visitor

type Member struct {
	ID        string
	Name      string
	ExternVal string
}

func (m Member) Accept(v MemberVisitor) error {
	return v.visitMember(m)
}

type Company struct {
	ID             string
	EntityName     string
	Classification string
}

type MemcoMember struct {
	ID        string
	FirstName string
	LastName  string
}

func (m MemcoMember) Accept(v MemberVisitor) error {
	return v.visitMemcoMember(m)
}

type PeoMember struct {
	ID        int
	FirstName string
	LastName  string
}

func (m PeoMember) Accept(v MemberVisitor) error {
	return v.visitPeoMember(m)
}

type EorMember struct {
	ID         int
	FirstName  string
	MiddleName string
	LastName   string
}

func (m EorMember) Accept(v MemberVisitor) error {
	return v.visitEorMember(m)
}
