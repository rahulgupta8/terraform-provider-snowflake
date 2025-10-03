//go:build !account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Contact_Basic(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	email := "test@example.com"
	updatedEmail := "updated@example.com"
	comment := random.Comment()

	contactModel := model.Contact("contact", id.Name(), email)
	contactModelWithComment := model.Contact("contact", id.Name(), email).
		WithComment(comment)
	contactModelWithUpdatedEmail := model.Contact("contact", id.Name(), updatedEmail).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			// create with empty optionals
			{
				Config: accconfig.FromModels(t, contactModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "email", email),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(contactModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.email", email),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.comment", ""),
				),
			},
			// import - without optionals
			{
				Config:       accconfig.FromModels(t, contactModel),
				ResourceName: contactModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "email", email),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set optionals
			{
				Config: accconfig.FromModels(t, contactModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(contactModelWithComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "email", email),
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(contactModelWithComment.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "show_output.0.email", email),
					resource.TestCheckResourceAttr(contactModelWithComment.ResourceReference(), "show_output.0.comment", comment),
				),
			},
			// import - complete
			{
				Config:       accconfig.FromModels(t, contactModelWithComment),
				ResourceName: contactModelWithComment.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: importchecks.ComposeAggregateImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "email", email),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// update email
			{
				Config: accconfig.FromModels(t, contactModelWithUpdatedEmail),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(contactModelWithUpdatedEmail.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "email", updatedEmail),
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "comment", comment),

					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(contactModelWithUpdatedEmail.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "show_output.0.email", updatedEmail),
					resource.TestCheckResourceAttr(contactModelWithUpdatedEmail.ResourceReference(), "show_output.0.comment", comment),
				),
			},
			// unset comment
			{
				Config: accconfig.FromModels(t, model.Contact("contact", id.Name(), updatedEmail)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(contactModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "email", updatedEmail),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "comment", ""),

					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.#", "1"),
					resource.TestCheckResourceAttrSet(contactModel.ResourceReference(), "show_output.0.created_on"),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.name", id.Name()),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.email", updatedEmail),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "show_output.0.comment", ""),
				),
			},
		},
	})
}

func TestAcc_Contact_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	email := "rename@example.com"

	contactModel := model.Contact("contact", id.Name(), email)
	contactModelRenamed := model.Contact("contact", newId.Name(), email)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.Contact),
		Steps: []resource.TestStep{
			// create
			{
				Config: accconfig.FromModels(t, contactModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "name", id.Name()),
					resource.TestCheckResourceAttr(contactModel.ResourceReference(), "email", email),
				),
			},
			// rename
			{
				Config: accconfig.FromModels(t, contactModelRenamed),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(contactModelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(contactModelRenamed.ResourceReference(), "name", newId.Name()),
					resource.TestCheckResourceAttr(contactModelRenamed.ResourceReference(), "email", email),
				),
			},
		},
	})
}
