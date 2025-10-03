package sdk

var (
	_ validatable = new(CreateContactOptions)
	_ validatable = new(AlterContactOptions)
	_ validatable = new(DropContactOptions)
	_ validatable = new(ShowContactOptions)
	_ validatable = new(DescribeContactOptions)
)

func (opts *CreateContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("CreateContactOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *AlterContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Set, opts.Unset, opts.RenameTo) {
		errs = append(errs, errExactlyOneOf("AlterContactOptions", "Set", "Unset", "RenameTo"))
	}
	if opts.RenameTo != nil && !ValidObjectIdentifier(opts.RenameTo) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if valueSet(opts.Set) {
		if !anyValueSet(opts.Set.Email, opts.Set.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterContactOptions.Set", "Email", "Comment"))
		}
	}
	if valueSet(opts.Unset) {
		if !anyValueSet(opts.Unset.Comment) {
			errs = append(errs, errAtLeastOneOf("AlterContactOptions.Unset", "Comment"))
		}
	}
	return JoinErrors(errs...)
}

func (opts *DropContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *ShowContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	return JoinErrors(errs...)
}

func (opts *DescribeContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}
