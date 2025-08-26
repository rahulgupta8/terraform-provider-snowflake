package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Contact_basic(t *testing.T) {
	contactName := random.AlphaN(8)
	email := fmt.Sprintf("test-%s@example.com", random.AlphaN(8))
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			{
				Config: contactConfig(contactName, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "fully_qualified_name"),
					resource.TestCheckResourceAttr("snowflake_contact.test", "show_output.#", "1"),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "show_output.0.created_on"),
					resource.TestCheckResourceAttr("snowflake_contact.test", "show_output.0.name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "show_output.0.email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "show_output.0.comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "show_output.0.owner"),
				),
			},
		},
	})
}

func TestAcc_Contact_update(t *testing.T) {
	contactName := random.AlphaN(8)
	email := fmt.Sprintf("test-%s@example.com", random.AlphaN(8))
	emailUpdated := fmt.Sprintf("updated-%s@example.org", random.AlphaN(8))
	comment := random.Comment()
	commentUpdated := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			// Create contact
			{
				Config: contactConfig(contactName, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "fully_qualified_name"),
				),
			},
			// Update email and comment
			{
				Config: contactConfig(contactName, emailUpdated, commentUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", emailUpdated),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", commentUpdated),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "fully_qualified_name"),
				),
			},
		},
	})
}

func TestAcc_Contact_rename(t *testing.T) {
	contactName := random.AlphaN(8)
	contactNameUpdated := random.AlphaN(8)
	email := fmt.Sprintf("test-%s@example.com", random.AlphaN(8))
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			// Create contact
			{
				Config: contactConfig(contactName, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
				),
			},
			// Rename contact
			{
				Config: contactConfig(contactNameUpdated, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactNameUpdated),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
				),
			},
		},
	})
}

func TestAcc_Contact_minimal(t *testing.T) {
	contactName := random.AlphaN(8)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			{
				Config: contactMinimalConfig(contactName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", ""),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", ""),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "fully_qualified_name"),
				),
			},
		},
	})
}

func TestAcc_Contact_unsetFields(t *testing.T) {
	contactName := random.AlphaN(8)
	email := fmt.Sprintf("test-%s@example.com", random.AlphaN(8))
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			// Create contact with email and comment
			{
				Config: contactConfig(contactName, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
				),
			},
			// Unset email and comment
			{
				Config: contactMinimalConfig(contactName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", ""),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", ""),
				),
			},
		},
	})
}

func TestAcc_Contact_import(t *testing.T) {
	contactName := random.AlphaN(8)
	email := fmt.Sprintf("test-%s@example.com", random.AlphaN(8))
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			{
				Config: contactConfig(contactName, email, comment),
			},
			{
				ResourceName:      "snowflake_contact.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func contactConfig(name, email, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_contact" "test" {
  name    = %[1]q
  email   = %[2]q
  comment = %[3]q
}
`, name, email, comment)
}

func contactMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_contact" "test" {
  name = %[1]q
}
`, name)
}
