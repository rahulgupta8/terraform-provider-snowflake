package sdk

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Contacts = (*contacts)(nil)
)

type contacts struct {
	client *Client
}

func (c *contacts) Create(ctx context.Context, request *CreateContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *contacts) Alter(ctx context.Context, request *AlterContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *contacts) Show(ctx context.Context, request *ShowContactRequest) ([]Contact, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[contactRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[contactRow, Contact](dbRows)
	return resultList, nil
}

func (c *contacts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Contact, error) {
	contacts, err := c.Show(ctx, NewShowContactRequest().WithLike(Like{Pattern: String(id.Name())}))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(contacts, func(r Contact) bool { return r.Name == id.Name() })
}

func (c *contacts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Contact, error) {
	contact, err := c.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, collections.ErrObjectNotFound) {
			return nil, ErrObjectNotFound
		}
		return nil, err
	}
	return contact, nil
}

func (c *contacts) Drop(ctx context.Context, request *DropContactRequest) error {
	opts := request.toOpts()
	return validateAndExec(c.client, ctx, opts)
}

func (c *contacts) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return c.Drop(ctx, NewDropContactRequest(id).WithIfExists(true))
}

func (r *CreateContactRequest) toOpts() *createContactOptions {
	opts := &createContactOptions{
		create:  true,
		contact: "CONTACT",
		name:    r.name,
	}
	if r.OrReplace != nil {
		opts.OrReplace = r.OrReplace
	}
	if r.IfNotExists != nil {
		opts.IfNotExists = r.IfNotExists
	}
	if r.Email != nil {
		opts.Email = r.Email
	}
	if r.Comment != nil {
		opts.Comment = r.Comment
	}
	return opts
}

func (r *AlterContactRequest) toOpts() *alterContactOptions {
	opts := &alterContactOptions{
		alter:   true,
		contact: "CONTACT",
		name:    r.name,
	}
	if r.IfExists != nil {
		opts.IfExists = r.IfExists
	}
	if r.Set != nil {
		opts.Set = r.Set
	}
	if r.Unset != nil {
		opts.Unset = r.Unset
	}
	if r.Rename != nil {
		opts.Rename = r.Rename
	}
	return opts
}

func (r *ShowContactRequest) toOpts() *showContactOptions {
	opts := &showContactOptions{
		show:    true,
		contact: "CONTACTS",
	}
	if r.Like != nil {
		opts.Like = r.Like
	}
	return opts
}

func (r *DropContactRequest) toOpts() *dropContactOptions {
	opts := &dropContactOptions{
		drop:    true,
		contact: "CONTACT",
		name:    r.name,
	}
	if r.IfExists != nil {
		opts.IfExists = r.IfExists
	}
	return opts
}

func (r contactRow) convert() *Contact {
	contact := &Contact{
		CreatedOn: r.CreatedOn,
		Name:      r.Name,
		Email:     r.Email,
		Comment:   r.Comment,
		Owner:     r.Owner,
	}
	return contact
}
