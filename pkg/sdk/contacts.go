package sdk

import (
	"context"
	"time"
)

type Contacts interface {
	Create(ctx context.Context, request *CreateContactRequest) error
	Alter(ctx context.Context, request *AlterContactRequest) error
	Show(ctx context.Context, opts *ShowContactRequest) ([]Contact, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Contact, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Contact, error)
	Drop(ctx context.Context, request *DropContactRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
}

type CreateContactRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        AccountObjectIdentifier
	Email       *string
	Comment     *string
}

type AlterContactRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier
	Set      *ContactSet
	Unset    *ContactUnset
	Rename   *ContactRename
}

type ShowContactRequest struct {
	Like *Like
}

type DropContactRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier
}

// createContactOptions is based on https://docs.snowflake.com/en/user-guide/contacts-using
type createContactOptions struct {
	create      bool                     `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	contact     string                   `ddl:"static" sql:"CONTACT"`
	IfNotExists *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier  `ddl:"identifier"`
	Email       *string                  `ddl:"parameter,single_quotes" sql:"EMAIL"`
	Comment     *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type alterContactOptions struct {
	alter       bool                     `ddl:"static" sql:"ALTER"`
	contact     string                   `ddl:"static" sql:"CONTACT"`
	IfExists    *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name        AccountObjectIdentifier  `ddl:"identifier"`
	Set         *ContactSet              `ddl:"keyword" sql:"SET"`
	Unset       *ContactUnset            `ddl:"keyword" sql:"UNSET"`
	Rename      *ContactRename           `ddl:"keyword" sql:"RENAME TO"`
}

type ContactSet struct {
	Email   *string `ddl:"parameter,single_quotes" sql:"EMAIL"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ContactUnset struct {
	Email   *bool `ddl:"keyword" sql:"EMAIL"`
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

type ContactRename struct {
	NewName AccountObjectIdentifier `ddl:"identifier"`
}

// showContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-contacts
type showContactOptions struct {
	show    bool   `ddl:"static" sql:"SHOW"`
	contact string `ddl:"static" sql:"CONTACTS"`
	Like    *Like  `ddl:"keyword" sql:"LIKE"`
}

// dropContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-contact
type dropContactOptions struct {
	drop     bool                     `ddl:"static" sql:"DROP"`
	contact  string                   `ddl:"static" sql:"CONTACT"`
	IfExists *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier  `ddl:"identifier"`
}

type Contact struct {
	CreatedOn time.Time
	Name      string
	Email     string
	Comment   string
	Owner     string
}

func (c *Contact) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(c.Name)
}

type contactRow struct {
	CreatedOn time.Time `db:"created_on"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Comment   string    `db:"comment"`
	Owner     string    `db:"owner"`
}