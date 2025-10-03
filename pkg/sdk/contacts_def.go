package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var ContactsDef = g.NewInterface(
	"Contacts",
	"Contact",
	g.KindOfT[AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-contact",
		g.NewQueryStruct("CreateContact").
			Create().
			OrReplace().
			SQL("NOTIFICATION CONTACT").
			IfNotExists().
			Name().
			TextAssignment("EMAIL", g.ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-contact",
		g.NewQueryStruct("AlterContact").
			Alter().
			SQL("NOTIFICATION CONTACT").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("ContactSet").
					OptionalTextAssignment("EMAIL", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "Email", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("ContactUnset").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RenameTo").
			WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-contact",
		g.NewQueryStruct("DropContact").
			Drop().
			SQL("NOTIFICATION CONTACT").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-contacts",
		g.DbStruct("showContactDBRow").
			Text("created_on").
			Text("name").
			Text("email").
			OptionalText("comment"),
		g.PlainStruct("Contact").
			Text("CreatedOn").
			Text("Name").
			Text("Email").
			Text("Comment"),
		g.NewQueryStruct("ShowContacts").
			Show().
			SQL("NOTIFICATION CONTACTS").
			OptionalLike(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-contact",
		g.DbStruct("describeContactDBRow").
			Text("property").
			Text("value"),
		g.PlainStruct("ContactProperty").
			Text("Name").
			Text("Value"),
		g.NewQueryStruct("DescribeContact").
			Describe().
			SQL("NOTIFICATION CONTACT").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
