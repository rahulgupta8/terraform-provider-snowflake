package testacc

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
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
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			{
				Config: contactConfig(contactName, email, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_contact.test", "name", contactName),
					resource.TestCheckResourceAttr("snowflake_contact.test", "email", email),
					resource.TestCheckResourceAttr("snowflake_contact.test", "comment", comment),
					resource.TestCheckResourceAttrSet("snowflake_contact.test", "fully_qualified_name"),
				),
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
