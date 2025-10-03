package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Contacts = (*contacts)(nil)

type contacts struct {
	client *Client
}

func (v *contacts) Create(ctx context.Context, request *CreateContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *contacts) Alter(ctx context.Context, request *AlterContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *contacts) Drop(ctx context.Context, request *DropContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *contacts) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropContactRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *contacts) Show(ctx context.Context, request *ShowContactRequest) ([]Contact, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[showContactDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[showContactDBRow, Contact](dbRows)
	return resultList, nil
}

func (v *contacts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Contact, error) {
	request := NewShowContactRequest().
		WithLike(Like{Pattern: String(id.Name())})
	contacts, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(contacts, func(r Contact) bool { return r.Name == id.Name() })
}

func (v *contacts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Contact, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *contacts) Describe(ctx context.Context, id AccountObjectIdentifier) ([]ContactProperty, error) {
	opts := &DescribeContactOptions{
		name: id,
	}
	rows, err := validateAndQuery[describeContactDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[describeContactDBRow, ContactProperty](rows), nil
}

func (r *CreateContactRequest) toOpts() *CreateContactOptions {
	opts := &CreateContactOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Email:       r.Email,
		Comment:     r.Comment,
	}
	return opts
}

func (r *AlterContactRequest) toOpts() *AlterContactOptions {
	opts := &AlterContactOptions{
		IfExists: r.IfExists,
		name:     r.name,

		RenameTo: r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &ContactSet{
			Email:   r.Set.Email,
			Comment: r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &ContactUnset{
			Comment: r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropContactRequest) toOpts() *DropContactOptions {
	opts := &DropContactOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowContactRequest) toOpts() *ShowContactOptions {
	opts := &ShowContactOptions{
		Like: r.Like,
	}
	return opts
}

func (r showContactDBRow) convert() *Contact {
	// TODO: Mapping
	return &Contact{}
}

func (r *DescribeContactRequest) toOpts() *DescribeContactOptions {
	opts := &DescribeContactOptions{
		name: r.name,
	}
	return opts
}

func (r describeContactDBRow) convert() *ContactProperty {
	// TODO: Mapping
	return &ContactProperty{}
}
