package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateContactOptions]   = new(CreateContactRequest)
	_ optionsProvider[AlterContactOptions]    = new(AlterContactRequest)
	_ optionsProvider[DropContactOptions]     = new(DropContactRequest)
	_ optionsProvider[ShowContactOptions]     = new(ShowContactRequest)
	_ optionsProvider[DescribeContactOptions] = new(DescribeContactRequest)
)

type CreateContactRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	Email       string                  // required
	Comment     *string
}

type AlterContactRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *ContactSetRequest
	Unset    *ContactUnsetRequest
	RenameTo *AccountObjectIdentifier
}

type ContactSetRequest struct {
	Email   *string
	Comment *string
}

type ContactUnsetRequest struct {
	Comment *bool
}

type DropContactRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowContactRequest struct {
	Like *Like
}

type DescribeContactRequest struct {
	name AccountObjectIdentifier // required
}
