package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowContactSchema represents output of SHOW query for the single Contact.
var ShowContactSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"email": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowContactSchema

func ContactToSchema(contact *sdk.Contact) map[string]any {
	contactSchema := make(map[string]any)
	contactSchema["created_on"] = contact.CreatedOn.String()
	contactSchema["name"] = contact.Name
	contactSchema["email"] = contact.Email
	contactSchema["comment"] = contact.Comment
	contactSchema["owner"] = contact.Owner
	return contactSchema
}
