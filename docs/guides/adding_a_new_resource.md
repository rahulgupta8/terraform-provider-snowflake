---
page_title: "Adding a New Resource - Contributor Guide"
subcategory: "Contributing"
---

# Adding a New Resource to the Snowflake Provider

This guide walks you through the complete process of adding a new resource to the Snowflake Terraform provider. It breaks down each step with examples and references to help you understand the structure and conventions used in the codebase.

## Prerequisites

Before you start adding a new resource, make sure you have:

1. **Set up your development environment** - Follow the instructions in [CONTRIBUTING.md](../../CONTRIBUTING.md#setting-up-the-development-environment)
2. **Discussed your change** - Open an issue or comment on an existing one to ensure the resource is aligned with the project's roadmap
3. **Reviewed the Snowflake documentation** - Understand the SQL commands and behavior of the Snowflake object you're implementing
4. **Configured your Snowflake test account** - Set up `~/.snowflake/config` with your test account credentials

## Overview of Steps

Adding a new resource involves these main steps:

1. [Generate SDK Interface and Implementation](#step-1-generate-sdk-interface-and-implementation)
2. [Create Resource Implementation](#step-2-create-resource-implementation)
3. [Register the Resource](#step-3-register-the-resource)
4. [Add Acceptance Tests](#step-4-add-acceptance-tests)
5. [Generate Documentation](#step-5-generate-documentation)
6. [Run Tests and Validation](#step-6-run-tests-and-validation)

## Step 1: Generate SDK Interface and Implementation

The SDK provides Go abstractions over Snowflake SQL commands. We use a generator to create most of the boilerplate code.

### 1.1 Create the SDK Definition File

Create a new file in `pkg/sdk/` with the name pattern `<object_name>_def.go`:

```go
//go:generate go run ./poc/main.go

package sdk

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc"

// Example: Database Role definition
var DatabaseRole = poc.NewInterface(
    "DatabaseRoles",
    "DatabaseRole",
    // ... define operations
)
```

**Key components to define:**

- **Show operation** - Lists objects
- **Describe operation** - Gets detailed information about an object
- **Create operation** - Creates a new object
- **Alter operation** - Modifies an existing object
- **Drop operation** - Deletes an object

For detailed examples, see:
- `pkg/sdk/poc/example/database_role_def.go` - Simple example
- `pkg/sdk/poc/README.md` - Complete generator documentation

### 1.2 Register Your Definition

Edit `pkg/sdk/poc/main.go` and add your definition to the `definitionMapping`:

```go
var definitionMapping = map[string]poc.Interface{
    // ... existing definitions
    "your_object_def.go": YourObject,
}
```

### 1.3 Generate SDK Files

Run the generator to create the SDK files:

```bash
make clean-generator-your_object run-generator-your_object
```

This will generate several files:
- `your_object_gen.go` - SDK interface and options structs
- `your_object_dto_gen.go` - Request DTOs
- `your_object_dto_builders_gen.go` - DTO builder methods
- `your_object_validations_gen.go` - Validation logic
- `your_object_impl_gen.go` - Implementation
- `your_object_gen_test.go` - Unit test placeholders

### 1.4 Add Integration Tests

Create integration tests in `pkg/sdk/testint/` to verify the SDK works with actual Snowflake API:

```go
//go:build !account_level_tests

package testint

import (
    "testing"
    
    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestInt_YourObjects(t *testing.T) {
    client := testClient(t)
    ctx := testContext(t)
    
    // Test Create
    id := sdk.NewDatabaseObjectIdentifier("TEST_DB", "TEST_OBJECT")
    err := client.YourObjects.Create(ctx, sdk.NewCreateYourObjectRequest(id))
    require.NoError(t, err)
    
    // Test Show
    objects, err := client.YourObjects.Show(ctx, nil)
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(objects), 1)
    
    // Test Describe
    details, err := client.YourObjects.Describe(ctx, id)
    require.NoError(t, err)
    assert.NotNil(t, details)
    
    // Test Drop
    err = client.YourObjects.Drop(ctx, sdk.NewDropYourObjectRequest(id))
    require.NoError(t, err)
}
```

Run integration tests:
```bash
make test-integration
```

## Step 2: Create Resource Implementation

Now create the Terraform resource that uses the SDK.

### 2.1 Define the Resource Schema

Create `pkg/resources/your_resource.go`:

```go
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
    "github.com/hashicorp/go-cty/cty"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var yourResourceSchema = map[string]*schema.Schema{
    "name": {
        Type:             schema.TypeString,
        Required:         true,
        Description:      "Specifies the identifier for the resource.",
        DiffSuppressFunc: suppressIdentifierQuoting,
    },
    "database": {
        Type:             schema.TypeString,
        Required:         true,
        ForceNew:         true,
        Description:      "The database in which to create the resource.",
        DiffSuppressFunc: suppressIdentifierQuoting,
    },
    "comment": {
        Type:        schema.TypeString,
        Optional:    true,
        Description: "Specifies a comment for the resource.",
    },
    ShowOutputAttributeName: {
        Type:        schema.TypeList,
        Computed:    true,
        Description: "Outputs the result of `SHOW` command.",
        Elem: &schema.Resource{
            Schema: schemas.ShowYourResourceSchema,
        },
    },
    FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}
```

**Schema Field Types:**

- `Required: true` - User must provide this value
- `Optional: true` - User can optionally provide this value
- `Computed: true` - Value is determined by Snowflake
- `ForceNew: true` - Changing this field requires recreating the resource
- `DiffSuppressFunc` - Custom logic to ignore certain differences

### 2.2 Implement CRUD Operations

```go
func YourResource() *schema.Resource {
    deleteFunc := ResourceDeleteContextFunc(
        sdk.ParseDatabaseObjectIdentifier,
        func(client *sdk.Client) DropSafelyFunc[sdk.DatabaseObjectIdentifier] {
            return client.YourObjects.DropSafely
        },
    )

    return &schema.Resource{
        SchemaVersion: 1,

        CreateContext: TrackingCreateWrapper(resources.YourResource, CreateYourResource),
        ReadContext:   TrackingReadWrapper(resources.YourResource, ReadYourResource),
        UpdateContext: TrackingUpdateWrapper(resources.YourResource, UpdateYourResource),
        DeleteContext: TrackingDeleteWrapper(resources.YourResource, deleteFunc),

        Description: "Resource used to manage your objects. For more information, check [documentation](https://docs.snowflake.com/...).",

        Schema: yourResourceSchema,
        Importer: &schema.ResourceImporter{
            StateContext: TrackingImportWrapper(resources.YourResource, ImportName[sdk.DatabaseObjectIdentifier]),
        },

        CustomizeDiff: TrackingCustomDiffWrapper(resources.YourResource, customdiff.All(
            ComputedIfAnyAttributeChanged(yourResourceSchema, ShowOutputAttributeName, "comment", "name"),
            ComputedIfAnyAttributeChanged(yourResourceSchema, FullyQualifiedNameAttributeName, "name"),
        )),

        StateUpgraders: []schema.StateUpgrader{
            {
                Version: 0,
                Type:    cty.EmptyObject,
                Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
            },
        },
        Timeouts: defaultTimeouts,
    }
}
```

**Create Function:**

```go
func CreateYourResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    client := meta.(*provider.Context).Client

    databaseName := d.Get("database").(string)
    name := d.Get("name").(string)
    id := sdk.NewDatabaseObjectIdentifier(databaseName, name)
    
    createRequest := sdk.NewCreateYourObjectRequest(id)

    if v, ok := d.GetOk("comment"); ok {
        createRequest.WithComment(v.(string))
    }

    err := client.YourObjects.Create(ctx, createRequest)
    if err != nil {
        return diag.FromErr(err)
    }

    d.SetId(helpers.EncodeResourceIdentifier(id))

    return ReadYourResource(ctx, d, meta)
}
```

**Read Function:**

```go
func ReadYourResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    client := meta.(*provider.Context).Client

    id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    object, err := client.YourObjects.ShowByIDSafely(ctx, id)
    if err != nil {
        if errors.Is(err, sdk.ErrObjectNotFound) {
            d.SetId("")
            return diag.Diagnostics{
                diag.Diagnostic{
                    Severity: diag.Warning,
                    Summary:  "Resource not found; marking it as removed",
                    Detail:   fmt.Sprintf("Resource id: %s, Err: %s", id.FullyQualifiedName(), err),
                },
            }
        }
        return diag.FromErr(err)
    }

    if err := d.Set("comment", object.Comment); err != nil {
        return diag.FromErr(err)
    }

    if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
        return diag.FromErr(err)
    }

    if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.YourResourceToSchema(object)}); err != nil {
        return diag.FromErr(err)
    }

    return nil
}
```

**Update Function:**

```go
func UpdateYourResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
    client := meta.(*provider.Context).Client

    id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }

    if d.HasChange("name") {
        newId := sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), d.Get("name").(string))

        err = client.YourObjects.Alter(ctx, sdk.NewAlterYourObjectRequest(id).WithRename(newId))
        if err != nil {
            return diag.FromErr(err)
        }

        d.SetId(helpers.EncodeResourceIdentifier(newId))
        id = newId
    }

    if d.HasChange("comment") {
        newComment := d.Get("comment").(string)
        err := client.YourObjects.Alter(ctx, sdk.NewAlterYourObjectRequest(id).WithSetComment(newComment))
        if err != nil {
            return diag.FromErr(err)
        }
    }

    return ReadYourResource(ctx, d, meta)
}
```

### 2.3 Create Show Output Schema

Create `pkg/schemas/your_resource.go` to define the show output schema:

```go
package schemas

import (
    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowYourResourceSchema = map[string]*schema.Schema{
    "name": {
        Type:     schema.TypeString,
        Computed: true,
    },
    "database_name": {
        Type:     schema.TypeString,
        Computed: true,
    },
    "comment": {
        Type:     schema.TypeString,
        Computed: true,
    },
    "created_on": {
        Type:     schema.TypeString,
        Computed: true,
    },
    // ... other fields from SHOW output
}

func YourResourceToSchema(resource *sdk.YourObject) map[string]any {
    return map[string]any{
        "name":          resource.Name,
        "database_name": resource.DatabaseName,
        "comment":       resource.Comment,
        "created_on":    resource.CreatedOn.String(),
        // ... map other fields
    }
}
```

## Step 3: Register the Resource

### 3.1 Add Resource Constant

Edit `pkg/provider/resources/resources.go` and add your resource constant:

```go
const (
    // ... existing resources
    YourResource resource = "snowflake_your_resource"
)
```

### 3.2 Register in Provider

Edit `pkg/provider/provider.go` and add your resource to the `getResources()` function:

```go
func getResources() map[string]*schema.Resource {
    return map[string]*schema.Resource{
        // ... existing resources
        "snowflake_your_resource": resources.YourResource(),
    }
}
```

### 3.3 Add to Preview Features (if needed)

If this is a preview feature, add it to `pkg/provider/previewfeatures/preview_features.go`:

```go
const (
    // ... existing features
    YourResourceResource feature = "snowflake_your_resource_resource"
)

var allPreviewFeatures = []feature{
    // ... existing features
    YourResourceResource,
}
```

Then wrap your resource read function with the preview feature check:

```go
ReadContext: TrackingReadWrapper(resources.YourResource, 
    PreviewFeatureReadWrapper(previewfeatures.YourResourceResource, ReadYourResource)),
```

## Step 4: Add Acceptance Tests

Acceptance tests verify the resource works end-to-end with actual Snowflake.

### 4.1 Create Acceptance Test File

Create `pkg/resources/your_resource_acceptance_test.go`:

```go
//go:build !account_level_tests

package resources_test

import (
    "fmt"
    "testing"

    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
    "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_YourResource_Basic(t *testing.T) {
    _ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

    id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
    comment := random.Comment()

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
        TerraformVersionChecks: []tfversion.TerraformVersionCheck{
            tfversion.RequireAbove(tfversion.Version1_5_0),
        },
        PreCheck:     func() { acc.TestAccPreCheck(t) },
        CheckDestroy: acc.CheckDestroy(t, resources.YourResource),
        Steps: []resource.TestStep{
            // Create with basic fields
            {
                Config: yourResourceConfig(id.DatabaseName(), id.Name(), comment),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("snowflake_your_resource.test", "name", id.Name()),
                    resource.TestCheckResourceAttr("snowflake_your_resource.test", "database", id.DatabaseName()),
                    resource.TestCheckResourceAttr("snowflake_your_resource.test", "comment", comment),
                ),
            },
            // Import
            {
                ResourceName:      "snowflake_your_resource.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
        },
    })
}

func yourResourceConfig(database, name, comment string) string {
    return fmt.Sprintf(`
resource "snowflake_your_resource" "test" {
  database = "%s"
  name     = "%s"
  comment  = "%s"
}
`, database, name, comment)
}
```

### 4.2 Test Different Scenarios

Create tests for:
- Basic creation
- Update operations
- Import functionality
- Edge cases (empty values, special characters, etc.)
- Error conditions

Example test for updates:

```go
func TestAcc_YourResource_Update(t *testing.T) {
    _ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

    id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
    comment1 := random.Comment()
    comment2 := random.Comment()

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
        TerraformVersionChecks: []tfversion.TerraformVersionCheck{
            tfversion.RequireAbove(tfversion.Version1_5_0),
        },
        PreCheck:     func() { acc.TestAccPreCheck(t) },
        CheckDestroy: acc.CheckDestroy(t, resources.YourResource),
        Steps: []resource.TestStep{
            {
                Config: yourResourceConfig(id.DatabaseName(), id.Name(), comment1),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("snowflake_your_resource.test", "comment", comment1),
                ),
            },
            {
                Config: yourResourceConfig(id.DatabaseName(), id.Name(), comment2),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("snowflake_your_resource.test", "comment", comment2),
                ),
            },
        },
    })
}
```

### 4.3 Run Acceptance Tests

```bash
# Run all acceptance tests for your resource
TF_ACC=1 go test -v ./pkg/resources -run TestAcc_YourResource

# Or use make command
make test-acceptance
```

## Step 5: Generate Documentation

Documentation is auto-generated from your code and templates.

### 5.1 Create Documentation Template

Create `templates/resources/your_resource.md.tmpl`:

```markdown
---
page_title: "{{.Type}} ({{.ProviderName}})"
subcategory: ""
description: |-
{{ .Description | plainmarkdown | trimspace | prefixlines "  " }}
---

# {{.Type}} ({{.ProviderName}})

{{ .Description | trimspace }}

## Example Usage

{{tffile "examples/resources/snowflake_your_resource/resource.tf"}}

{{ .SchemaMarkdown | trimspace }}

## Import

Import is supported using the following syntax:

{{codefile "shell" "examples/resources/snowflake_your_resource/import.sh"}}
```

### 5.2 Create Example Files

Create `examples/resources/snowflake_your_resource/resource.tf`:

```hcl
# Basic resource
resource "snowflake_your_resource" "example" {
  database = "EXAMPLE_DB"
  name     = "example_resource"
  comment  = "This is an example resource"
}

# Resource with all options
resource "snowflake_your_resource" "complete" {
  database = "EXAMPLE_DB"
  name     = "complete_example"
  comment  = "Resource with all options set"
  
  # Add other optional parameters
}
```

Create `examples/resources/snowflake_your_resource/import.sh`:

```bash
terraform import snowflake_your_resource.example "database_name|resource_name"
```

### 5.3 Generate Documentation

```bash
make docs
```

This will generate `docs/resources/your_resource.md` with the complete documentation.

## Step 6: Run Tests and Validation

Before submitting your PR, run all validation steps:

### 6.1 Format Code

```bash
make fmt
```

### 6.2 Run Linter

```bash
make lint
```

### 6.3 Run Unit Tests

```bash
make test-unit
```

### 6.4 Run Integration Tests

```bash
make test-integration
```

### 6.5 Run Acceptance Tests

```bash
make test-acceptance
```

### 6.6 Complete Pre-Push Check

```bash
make pre-push
```

This runs all checks: formatting, linting, documentation generation, and tests.

## Best Practices and Tips

### Code Style

1. **Follow existing patterns** - Look at similar resources for inspiration
2. **Use consistent naming** - Follow the naming conventions in the codebase
3. **Add helpful error messages** - Make errors clear and actionable
4. **Document complex logic** - Add comments for non-obvious behavior

### Testing

1. **Test all CRUD operations** - Create, read, update, and delete
2. **Test edge cases** - Empty values, special characters, large inputs
3. **Test error conditions** - Invalid inputs, permission errors
4. **Use random values** - Avoid hardcoded values that might conflict
5. **Clean up resources** - Ensure tests clean up after themselves

### Documentation

1. **Provide clear examples** - Show common use cases
2. **Document all parameters** - Explain what each field does
3. **Link to Snowflake docs** - Reference official documentation
4. **Include import examples** - Show how to import existing resources

### Common Pitfalls

1. **Identifier handling** - Make sure to properly handle quoted identifiers
2. **State management** - Properly handle computed fields and diffs
3. **Error handling** - Don't mask errors, return them with context
4. **Preview features** - Remember to gate preview features appropriately

## Reference Resources

### Example Resources to Study

**Simple resources:**
- `pkg/resources/database_role.go` - Basic CRUD operations
- `pkg/resources/database.go` - Resource with parameters

**Complex resources:**
- `pkg/resources/schema.go` - Handling special cases (PUBLIC schema)
- `pkg/resources/user.go` - Multiple update operations

### Key Files to Reference

- `pkg/sdk/poc/README.md` - SDK generator documentation
- `pkg/resources/common.go` - Common helper functions
- `pkg/schemas/` - Show output schema examples
- `CONTRIBUTING.md` - General contribution guidelines
- `MIGRATION_GUIDE.md` - For breaking changes

## Getting Help

If you get stuck:

1. **Check existing resources** - Find similar resources and study their implementation
2. **Read the FAQ** - Check [FAQ.md](../../FAQ.md) for common questions
3. **Ask in discussions** - Use GitHub Discussions for questions
4. **Open an issue** - If you think you found a bug or missing feature
5. **Review the roadmap** - Check [ROADMAP.md](../../ROADMAP.md) for planned changes

## Submitting Your Contribution

Once you've completed all steps:

1. **Run pre-push checks** - `make pre-push` should pass
2. **Commit your changes** - Use [Conventional Commits](https://www.conventionalcommits.org/) format
3. **Push your branch** - Push to your fork
4. **Open a Pull Request** - Reference the issue you're addressing
5. **Wait for review** - Maintainers typically respond within 1-2 days

Your PR description should include:
- Link to the related issue
- Summary of changes
- Summary of tests added
- Any breaking changes or migration notes

Thank you for contributing to the Snowflake Terraform Provider! 🎉
