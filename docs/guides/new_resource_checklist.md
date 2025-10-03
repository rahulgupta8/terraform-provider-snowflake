---
page_title: "New Resource Checklist - Quick Reference"
subcategory: "Contributing"
---

# New Resource Implementation Checklist

This is a quick reference checklist for adding a new resource. For detailed instructions, see the [comprehensive guide](adding_a_new_resource.md).

## Before You Start

- [ ] Created/commented on a GitHub issue to discuss the change
- [ ] Reviewed Snowflake documentation for the object you're implementing
- [ ] Set up development environment (`make dev-setup`)
- [ ] Configured `~/.snowflake/config` with test credentials

## SDK Layer

### SDK Definition and Generation

- [ ] Created `pkg/sdk/your_object_def.go` with complete interface definition
  - [ ] Defined Show operation
  - [ ] Defined Describe operation (if applicable)
  - [ ] Defined Create operation
  - [ ] Defined Alter operations
  - [ ] Defined Drop operation
- [ ] Added definition to `definitionMapping` in `pkg/sdk/poc/main.go`
- [ ] Generated SDK files: `make clean-generator-your_object run-generator-your_object`
- [ ] Reviewed generated files and filled in any TODOs

### SDK Integration Tests

- [ ] Created `pkg/sdk/testint/your_objects_integration_test.go`
- [ ] Added tests for Create operation
- [ ] Added tests for Show operation
- [ ] Added tests for ShowByID operation
- [ ] Added tests for Describe operation (if applicable)
- [ ] Added tests for Alter operations
- [ ] Added tests for Drop operation
- [ ] All integration tests pass: `make test-integration`

## Resource Layer

### Resource Implementation

- [ ] Created `pkg/resources/your_resource.go` with:
  - [ ] Schema definition (`yourResourceSchema`)
  - [ ] Resource function (`YourResource()`)
  - [ ] Create function (`CreateYourResource`)
  - [ ] Read function (`ReadYourResource`)
  - [ ] Update function (`UpdateYourResource`)
  - [ ] Proper error handling in all functions
  - [ ] CustomizeDiff for computed fields
  - [ ] State upgraders (if needed)
  - [ ] Importer configuration
  - [ ] Description and documentation link

### Show Output Schema

- [ ] Created `pkg/schemas/your_resource.go` with:
  - [ ] `ShowYourResourceSchema` map
  - [ ] `YourResourceToSchema()` conversion function
  - [ ] All relevant fields from Snowflake SHOW command

### Resource Registration

- [ ] Added resource constant to `pkg/provider/resources/resources.go`
- [ ] Registered resource in `pkg/provider/provider.go` `getResources()` function
- [ ] If preview feature: Added to `pkg/provider/previewfeatures/preview_features.go`
- [ ] If preview feature: Wrapped Read function with `PreviewFeatureReadWrapper`

## Testing

### Acceptance Tests

- [ ] Created `pkg/resources/your_resource_acceptance_test.go` with:
  - [ ] Basic create test
  - [ ] Update test
  - [ ] Import test
  - [ ] Multiple configuration variations
  - [ ] Edge cases (empty values, special characters)
- [ ] All acceptance tests pass: `TF_ACC=1 go test -v ./pkg/resources -run TestAcc_YourResource`

## Documentation

### Template and Examples

- [ ] Created `templates/resources/your_resource.md.tmpl`
- [ ] Created `examples/resources/snowflake_your_resource/resource.tf` with:
  - [ ] Basic example
  - [ ] Complete example with all options
- [ ] Created `examples/resources/snowflake_your_resource/import.sh`
- [ ] Generated documentation: `make docs`
- [ ] Verified generated docs in `docs/resources/your_resource.md`

## Validation

### Code Quality

- [ ] Code is formatted: `make fmt`
- [ ] Linter passes: `make lint`
- [ ] No obvious code duplication
- [ ] Error messages are clear and helpful
- [ ] Comments added for complex logic

### Testing

- [ ] Unit tests pass: `make test-unit`
- [ ] Integration tests pass: `make test-integration`
- [ ] Acceptance tests pass: `make test-acceptance`
- [ ] Architecture tests pass: `make test-architecture`

### Final Check

- [ ] All pre-push checks pass: `make pre-push`
- [ ] Reviewed all changed files
- [ ] No debug code or console logs left
- [ ] Documentation is clear and accurate

## Pull Request

### PR Preparation

- [ ] Used Conventional Commits format for commit messages
- [ ] PR title follows Conventional Commits format
- [ ] PR description includes:
  - [ ] Link to related issue
  - [ ] Summary of changes
  - [ ] Summary of tests added
  - [ ] Any breaking changes or migration notes
- [ ] No unrelated changes included
- [ ] Commits are logical and well-organized

### After Submission

- [ ] Responded to review comments
- [ ] Made requested changes
- [ ] Re-ran tests after changes
- [ ] Marked resolved conversations

## Common Files to Review

When implementing your resource, these files can serve as references:

| Purpose | Simple Example | Complex Example |
|---------|---------------|-----------------|
| Resource Implementation | `pkg/resources/database_role.go` | `pkg/resources/schema.go` |
| SDK Definition | `pkg/sdk/poc/example/database_role_def.go` | `pkg/sdk/warehouses_def.go` |
| Integration Tests | `pkg/sdk/testint/database_roles_integration_test.go` | `pkg/sdk/testint/warehouses_integration_test.go` |
| Acceptance Tests | `pkg/resources/database_role_acceptance_test.go` | `pkg/resources/warehouse_acceptance_test.go` |
| Show Schema | `pkg/schemas/database_role.go` | `pkg/schemas/warehouse.go` |

## Quick Command Reference

```bash
# Setup
make dev-setup

# Generate SDK
make clean-generator-your_object run-generator-your_object

# Format code
make fmt

# Lint
make lint

# Test
make test-unit
make test-integration
make test-acceptance
TF_ACC=1 go test -v ./pkg/resources -run TestAcc_YourResource

# Generate docs
make docs

# Complete check
make pre-push
```

## Need Help?

- **Detailed guide**: [Adding a New Resource Guide](adding_a_new_resource.md)
- **General contribution**: [CONTRIBUTING.md](../../CONTRIBUTING.md)
- **FAQ**: [FAQ.md](../../FAQ.md)
- **GitHub Discussions**: Ask questions in the discussions area
- **Issues**: Report bugs or request features

---

**Pro Tip**: Don't try to do everything at once. Work incrementally:
1. Get SDK working first
2. Then create basic resource (create/read/delete)
3. Add update operations
4. Add tests
5. Generate documentation
6. Run all validation checks
