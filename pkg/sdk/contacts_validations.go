package sdk

var (
	_ validatable = new(createContactOptions)
	_ validatable = new(alterContactOptions)
	_ validatable = new(dropContactOptions)
	_ validatable = new(showContactOptions)
)

func (opts *createContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.IfNotExists, opts.OrReplace) {
		errs = append(errs, errOneOf("createContactOptions", "IfNotExists", "OrReplace"))
	}
	return JoinErrors(errs...)
}

func (opts *alterContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !anyValueSet(opts.Set, opts.Unset, opts.Rename) {
		errs = append(errs, errAtLeastOneOf("alterContactOptions", "Set", "Unset", "Rename"))
	}
	return JoinErrors(errs...)
}

func (opts *dropContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return JoinErrors(errs...)
}

func (opts *showContactOptions) validate() error {
	if opts == nil {
		return ErrNilOptions
	}
	return nil
}