package sdk

func NewCreateContactRequest(
	name AccountObjectIdentifier,
) *CreateContactRequest {
	s := CreateContactRequest{}
	s.name = name
	return &s
}

func (s *CreateContactRequest) WithOrReplace(orReplace bool) *CreateContactRequest {
	s.OrReplace = &orReplace
	return s
}

func (s *CreateContactRequest) WithIfNotExists(ifNotExists bool) *CreateContactRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateContactRequest) WithEmail(email string) *CreateContactRequest {
	s.Email = &email
	return s
}

func (s *CreateContactRequest) WithComment(comment string) *CreateContactRequest {
	s.Comment = &comment
	return s
}

func NewAlterContactRequest(
	name AccountObjectIdentifier,
) *AlterContactRequest {
	s := AlterContactRequest{}
	s.name = name
	return &s
}

func (s *AlterContactRequest) WithIfExists(ifExists bool) *AlterContactRequest {
	s.IfExists = &ifExists
	return s
}

func (s *AlterContactRequest) WithSet(set ContactSet) *AlterContactRequest {
	s.Set = &set
	return s
}

func (s *AlterContactRequest) WithUnset(unset ContactUnset) *AlterContactRequest {
	s.Unset = &unset
	return s
}

func (s *AlterContactRequest) WithRename(rename ContactRename) *AlterContactRequest {
	s.Rename = &rename
	return s
}

func NewContactSetRequest() *ContactSet {
	return &ContactSet{}
}

func (s *ContactSet) WithEmail(email string) *ContactSet {
	s.Email = &email
	return s
}

func (s *ContactSet) WithComment(comment string) *ContactSet {
	s.Comment = &comment
	return s
}

func NewContactUnsetRequest() *ContactUnset {
	return &ContactUnset{}
}

func (s *ContactUnset) WithEmail(email bool) *ContactUnset {
	s.Email = &email
	return s
}

func (s *ContactUnset) WithComment(comment bool) *ContactUnset {
	s.Comment = &comment
	return s
}

func NewShowContactRequest() *ShowContactRequest {
	return &ShowContactRequest{}
}

func (s *ShowContactRequest) WithLike(like Like) *ShowContactRequest {
	s.Like = &like
	return s
}

func NewDropContactRequest(
	name AccountObjectIdentifier,
) *DropContactRequest {
	s := DropContactRequest{}
	s.name = name
	return &s
}

func (s *DropContactRequest) WithIfExists(ifExists bool) *DropContactRequest {
	s.IfExists = &ifExists
	return s
}
