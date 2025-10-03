package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var contactSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("Identifier for the notification contact; must be unique for your account."),
	},
	"email": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Email address for the notification contact.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comment for the notification contact.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW NOTIFICATION CONTACTS` for the given contact.",
		Elem: &schema.Resource{
			Schema: schemas.ShowContactSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Contact() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] { return client.Contacts.DropSafely },
	)

	return &schema.Resource{
		Schema: contactSchema,

		CreateContext: TrackingCreateWrapper(resources.Contact, CreateContact),
		ReadContext:   TrackingReadWrapper(resources.Contact, ReadContact),
		DeleteContext: TrackingDeleteWrapper(resources.Contact, deleteFunc),
		UpdateContext: TrackingUpdateWrapper(resources.Contact, UpdateContact),
		Description:   "The resource is used for notification contact management. Notification contacts are email addresses used for alerts and notifications. For more details, refer to the [official documentation](https://docs.snowflake.com/en/sql-reference/sql/create-contact).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Contact, customdiff.All(
			ComputedIfAnyAttributeChanged(contactSchema, ShowOutputAttributeName, "comment", "name", "email"),
			ComputedIfAnyAttributeChanged(contactSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Contact, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	email := d.Get("email").(string)
	req := sdk.NewCreateContactRequest(id, email)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	err = client.Contacts.Create(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create notification contact",
				Detail:   fmt.Sprintf("Contact name: %s, err: %s", id.Name(), err),
			},
		}
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContact(ctx, d, meta)
}

func ReadContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	contact, err := client.Contacts.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Notification contact not found; marking it as removed",
					Detail:   fmt.Sprintf("Contact name: %s, err: %s", id.Name(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show notification contact by id",
				Detail:   fmt.Sprintf("Contact name: %s, err: %s", id.Name(), err),
			},
		}
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email", contact.Email); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set notification contact email",
				Detail:   fmt.Sprintf("Contact name: %s, email: %s, err: %s", contact.ID().FullyQualifiedName(), contact.Email, err),
			},
		}
	}

	if err := d.Set("comment", contact.Comment); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set notification contact comment",
				Detail:   fmt.Sprintf("Contact name: %s, comment: %s, err: %s", contact.ID().FullyQualifiedName(), contact.Comment, err),
			},
		}
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.ContactToSchema(contact)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithRenameTo(newId)); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to rename notification contact",
					Detail:   fmt.Sprintf("Previous contact name: %s, new contact name: %s, err: %s", id.Name(), newId.Name(), err),
				},
			}
		}

		id = newId
		d.SetId(helpers.EncodeResourceIdentifier(newId))
	}

	// Handle set operations
	if d.HasChange("email") || d.HasChange("comment") {
		set := sdk.NewContactSetRequest()
		if d.HasChange("email") {
			set.WithEmail(d.Get("email").(string))
		}
		if d.HasChange("comment") {
			if v, ok := d.GetOk("comment"); ok {
				set.WithComment(v.(string))
			} else {
				// If comment was removed, we need to unset it
				if err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithUnset(*sdk.NewContactUnsetRequest().WithComment(true))); err != nil {
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Failed to unset notification contact comment",
							Detail:   fmt.Sprintf("Contact name: %s, err: %s", id.Name(), err),
						},
					}
				}
				// Read after unset to get the updated state
				return ReadContact(ctx, d, meta)
			}
		}
		if err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithSet(*set)); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to update notification contact",
					Detail:   fmt.Sprintf("Contact name: %s, err: %s", id.Name(), err),
				},
			}
		}
	}

	return ReadContact(ctx, d, meta)
}
