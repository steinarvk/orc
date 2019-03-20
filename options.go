package orc

import "github.com/spf13/pflag"

type orcOptions struct {
	flags *pflag.FlagSet
}

func (o *orcOptions) intoUseContext(uc *UseContext) {
	uc.Flags = o.flags
}

func parseOptions(options []Option) (*orcOptions, error) {
	rv := &orcOptions{}
	for _, option := range options {
		if err := option(rv); err != nil {
			return nil, err
		}
	}
	return rv, nil
}

type Option func(*orcOptions) error

func WithFlags(flagset *pflag.FlagSet) Option {
	return func(opts *orcOptions) error {
		opts.flags = flagset
		return nil
	}
}
