package sdk

import "testing"

func TestContacts_Create(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid CreateContactOptions
	defaultOpts := func() *CreateContactOptions {
		return &CreateContactOptions{
			name:  id,
			Email: "test@example.com",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateContactOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateContactOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestContacts_Alter(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid AlterContactOptions
	defaultOpts := func() *AlterContactOptions {
		return &AlterContactOptions{
			name: id,
			Set:  &ContactSet{Email: String("new@example.com")},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterContactOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = nil
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterContactOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := &AlterContactOptions{
			name:     id,
			RenameTo: &emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of the fields [opts.Set.Email opts.Set.Comment] should be set", func(t *testing.T) {
		opts := &AlterContactOptions{
			name: id,
			Set:  &ContactSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterContactOptions.Set", "Email", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := &AlterContactOptions{
			name:  id,
			Unset: &ContactUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterContactOptions.Unset", "Comment"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `ALTER NOTIFICATION CONTACT %s SET EMAIL = 'new@example.com'`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &AlterContactOptions{
			name:     id,
			IfExists: Bool(true),
			Set: &ContactSet{
				Email:   String("updated@example.com"),
				Comment: String("updated comment"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER NOTIFICATION CONTACT IF EXISTS %s SET EMAIL = 'updated@example.com' COMMENT = 'updated comment'`, id.FullyQualifiedName())
	})
}

func TestContacts_Drop(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid DropContactOptions
	defaultOpts := func() *DropContactOptions {
		return &DropContactOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropContactOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP NOTIFICATION CONTACT %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP NOTIFICATION CONTACT IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestContacts_Show(t *testing.T) {
	// Minimal valid ShowContactOptions
	defaultOpts := func() *ShowContactOptions {
		return &ShowContactOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowContactOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW NOTIFICATION CONTACTS`)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("test_contact")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW NOTIFICATION CONTACTS LIKE 'test_contact'`)
	})
}

func TestContacts_Describe(t *testing.T) {

	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeContactOptions
	defaultOpts := func() *DescribeContactOptions {
		return &DescribeContactOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeContactOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE NOTIFICATION CONTACT %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE NOTIFICATION CONTACT %s`, id.FullyQualifiedName())
	})
}
