package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
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
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the contact; must be unique for the account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"email": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Email address for the contact.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the contact.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CONTACTS` for the given contact.",
		Elem: &schema.Resource{
			Schema: schemas.ShowContactSchema,
		},
	},
}

func Contact() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.Contact, CreateContextContact),
		ReadContext:   TrackingReadWrapper(resources.Contact, ReadContextContact),
		UpdateContext: TrackingUpdateWrapper(resources.Contact, UpdateContextContact),
		DeleteContext: TrackingDeleteWrapper(resources.Contact, DeleteContextContact),
		Description:   "Resource used to manage contact objects. For more information, check [contact documentation](https://docs.snowflake.com/en/user-guide/contacts-using).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Contact, customdiff.All(
			ComputedIfAnyAttributeChanged(contactSchema, ShowOutputAttributeName, "name", "email", "comment"),
			ComputedIfAnyAttributeChanged(contactSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: contactSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Contact, ImportName[sdk.AccountObjectIdentifier]),
		},
	}
}

func CreateContextContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)

	request := sdk.NewCreateContactRequest(id)
	if v, ok := d.GetOk("email"); ok {
		request.WithEmail(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.Contacts.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextContact(ctx, d, meta)
}

func ReadContextContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query contact. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Contact id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ContactToSchema(contact)}),
		d.Set("name", contact.Name),
		d.Set("email", contact.Email),
		d.Set("comment", contact.Comment),
	); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithRename(sdk.ContactRename{NewName: newId}))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming contact %v err = %w", d.Id(), err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	var sets []interface{}
	var unsets []interface{}

	if d.HasChange("email") {
		email, ok := d.GetOk("email")
		if ok {
			sets = append(sets, sdk.NewContactSetRequest().WithEmail(email.(string)))
		} else {
			unsets = append(unsets, sdk.NewContactUnsetRequest().WithEmail(true))
		}
	}

	if d.HasChange("comment") {
		comment, ok := d.GetOk("comment")
		if ok {
			sets = append(sets, sdk.NewContactSetRequest().WithComment(comment.(string)))
		} else {
			unsets = append(unsets, sdk.NewContactUnsetRequest().WithComment(true))
		}
	}

	for _, set := range sets {
		if err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithSet(*set.(*sdk.ContactSet))); err != nil {
			return diag.FromErr(err)
		}
	}

	for _, unset := range unsets {
		if err := client.Contacts.Alter(ctx, sdk.NewAlterContactRequest(id).WithUnset(*unset.(*sdk.ContactUnset))); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextContact(ctx, d, meta)
}

func DeleteContextContact(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Contacts.Drop(ctx, sdk.NewDropContactRequest(id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
