//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Contacts(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertContact := func(t *testing.T, contact *sdk.Contact, id sdk.AccountObjectIdentifier, expectedEmail string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, contact.CreatedOn)
		assert.Equal(t, id.Name(), contact.Name)
		assert.Equal(t, expectedEmail, contact.Email)
		assert.Equal(t, expectedComment, contact.Comment)
	}

	cleanupContactProvider := func(id sdk.AccountObjectIdentifier) func() {
		return func() {
			err := client.Contacts.Drop(ctx, sdk.NewDropContactRequest(id))
			require.NoError(t, err)
		}
	}

	createContact := func(t *testing.T, email string) *sdk.Contact {
		t.Helper()
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Contacts.Create(ctx, sdk.NewCreateContactRequest(id, email))
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(id))

		contact, err := client.Contacts.ShowByID(ctx, id)
		require.NoError(t, err)

		return contact
	}

	t.Run("create contact: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		email := "test@example.com"
		comment := random.Comment()

		request := sdk.NewCreateContactRequest(id, email).
			WithComment(&comment).
			WithIfNotExists(sdk.Bool(true))

		err := client.Contacts.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(id))

		contact, err := client.Contacts.ShowByID(ctx, id)

		require.NoError(t, err)
		assertContact(t, contact, id, email, comment)
	})

	t.Run("create contact: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		email := "test2@example.com"

		request := sdk.NewCreateContactRequest(id, email)

		err := client.Contacts.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(id))

		contact, err := client.Contacts.ShowByID(ctx, id)

		require.NoError(t, err)
		assertContact(t, contact, id, email, "")
	})

	t.Run("drop contact: existing", func(t *testing.T) {
		email := "drop@example.com"
		contact := createContact(t, email)

		err := client.Contacts.Drop(ctx, sdk.NewDropContactRequest(contact.ID()))
		require.NoError(t, err)

		_, err = client.Contacts.ShowByID(ctx, contact.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop contact: non-existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Contacts.Drop(ctx, sdk.NewDropContactRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	})

	t.Run("alter contact: set email", func(t *testing.T) {
		email := "original@example.com"
		contact := createContact(t, email)

		newEmail := "updated@example.com"
		alterRequest := sdk.NewAlterContactRequest(contact.ID()).WithSet(
			sdk.NewContactSetRequest().WithEmail(&newEmail),
		)
		err := client.Contacts.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedContact, err := client.Contacts.ShowByID(ctx, contact.ID())
		require.NoError(t, err)
		assert.Equal(t, newEmail, updatedContact.Email)
	})

	t.Run("alter contact: set comment", func(t *testing.T) {
		email := "comment@example.com"
		contact := createContact(t, email)

		newComment := random.Comment()
		alterRequest := sdk.NewAlterContactRequest(contact.ID()).WithSet(
			sdk.NewContactSetRequest().WithComment(&newComment),
		)
		err := client.Contacts.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedContact, err := client.Contacts.ShowByID(ctx, contact.ID())
		require.NoError(t, err)
		assert.Equal(t, newComment, updatedContact.Comment)
	})

	t.Run("alter contact: unset comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		email := "unset@example.com"
		comment := random.Comment()

		request := sdk.NewCreateContactRequest(id, email).WithComment(&comment)
		err := client.Contacts.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(id))

		alterRequest := sdk.NewAlterContactRequest(id).WithUnset(
			sdk.NewContactUnsetRequest().WithComment(sdk.Bool(true)),
		)
		err = client.Contacts.Alter(ctx, alterRequest)
		require.NoError(t, err)

		updatedContact, err := client.Contacts.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", updatedContact.Comment)
	})

	t.Run("alter contact: rename", func(t *testing.T) {
		email := "rename@example.com"
		contact := createContact(t, email)

		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		alterRequest := sdk.NewAlterContactRequest(contact.ID()).WithRenameTo(&newId)
		err := client.Contacts.Alter(ctx, alterRequest)
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(newId))

		_, err = client.Contacts.ShowByID(ctx, contact.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

		renamedContact, err := client.Contacts.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, newId.Name(), renamedContact.Name)
		assert.Equal(t, email, renamedContact.Email)
	})

	t.Run("show contact: without like", func(t *testing.T) {
		email1 := "show1@example.com"
		email2 := "show2@example.com"
		contact1 := createContact(t, email1)
		contact2 := createContact(t, email2)

		contacts, err := client.Contacts.Show(ctx, sdk.NewShowContactRequest())
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(contacts), 2)
		assert.Contains(t, contacts, *contact1)
		assert.Contains(t, contacts, *contact2)
	})

	t.Run("show contact: with like", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		email := "like@example.com"

		request := sdk.NewCreateContactRequest(id, email)
		err := client.Contacts.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupContactProvider(id))

		contacts, err := client.Contacts.Show(ctx, sdk.NewShowContactRequest().WithLike(&sdk.Like{Pattern: &id.name}))
		require.NoError(t, err)

		assert.Len(t, contacts, 1)
		assert.Equal(t, id.Name(), contacts[0].Name)
	})

	t.Run("describe contact", func(t *testing.T) {
		email := "describe@example.com"
		contact := createContact(t, email)

		properties, err := client.Contacts.Describe(ctx, contact.ID())
		require.NoError(t, err)
		assert.NotEmpty(t, properties)
	})
}
