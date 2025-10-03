package sdk

import (
	"context"
	"database/sql"
)

type Contacts interface {
	Create(ctx context.Context, request *CreateContactRequest) error
	Alter(ctx context.Context, request *AlterContactRequest) error
	Drop(ctx context.Context, request *DropContactRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowContactRequest) ([]Contact, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Contact, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Contact, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ContactProperty, error)
}

// CreateContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-contact.
type CreateContactOptions struct {
	create              bool                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	notificationContact bool                    `ddl:"static" sql:"NOTIFICATION CONTACT"`
	IfNotExists         *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
	Email               string                  `ddl:"parameter,single_quotes" sql:"EMAIL"`
	Comment             *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-contact.
type AlterContactOptions struct {
	alter               bool                     `ddl:"static" sql:"ALTER"`
	notificationContact bool                     `ddl:"static" sql:"NOTIFICATION CONTACT"`
	IfExists            *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier  `ddl:"identifier"`
	Set                 *ContactSet              `ddl:"keyword" sql:"SET"`
	Unset               *ContactUnset            `ddl:"list,no_parentheses" sql:"UNSET"`
	RenameTo            *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
}

type ContactSet struct {
	Email   *string `ddl:"parameter,single_quotes" sql:"EMAIL"`
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ContactUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-contact.
type DropContactOptions struct {
	drop                bool                    `ddl:"static" sql:"DROP"`
	notificationContact bool                    `ddl:"static" sql:"NOTIFICATION CONTACT"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                AccountObjectIdentifier `ddl:"identifier"`
}

// ShowContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-contacts.
type ShowContactOptions struct {
	show                 bool  `ddl:"static" sql:"SHOW"`
	notificationContacts bool  `ddl:"static" sql:"NOTIFICATION CONTACTS"`
	Like                 *Like `ddl:"keyword" sql:"LIKE"`
}

type showContactDBRow struct {
	CreatedOn string         `db:"created_on"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	Comment   sql.NullString `db:"comment"`
}

type Contact struct {
	CreatedOn string
	Name      string
	Email     string
	Comment   string
}

func (v *Contact) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
func (v *Contact) ObjectType() ObjectType {
	return ObjectTypeContact
}

// DescribeContactOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-contact.
type DescribeContactOptions struct {
	describe            bool                    `ddl:"static" sql:"DESCRIBE"`
	notificationContact bool                    `ddl:"static" sql:"NOTIFICATION CONTACT"`
	name                AccountObjectIdentifier `ddl:"identifier"`
}

type describeContactDBRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type ContactProperty struct {
	Name  string
	Value string
}
