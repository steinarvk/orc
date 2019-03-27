package orc

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func mkPositionalArgsError(args []string) error {
	return fmt.Errorf("expected no positional arguments, got: %v", args)
}

func castBodyFunction(body interface{}) (func([]string) error, error) {
	switch f := body.(type) {
	case func([]string) error:
		return f, nil
	case func():
		return func(args []string) error {
			if len(args) > 0 {
				return mkPositionalArgsError(args)
			}
			f()
			return nil
		}, nil
	case func([]string):
		return func(args []string) error {
			f(args)
			return nil
		}, nil
	case func() error:
		return func(args []string) error {
			if len(args) > 0 {
				return mkPositionalArgsError(args)
			}
			return f()
		}, nil
	default:
		return nil, errors.New("unsupported Body function supplied")
	}
}

func ReuseCommand(parent *cobra.Command, existing *cobra.Command) *cobra.Command {
	cmd := *existing
	parent.AddCommand(&cmd)
	return &cmd
}

func Command(parent *cobra.Command, prereq Module, skeleton cobra.Command, body interface{}, options ...Option) *cobra.Command {
	cmd := skeleton

	options = append(options, WithFlags(cmd.Flags()))

	if prereq == nil {
		prereq = Modules()
	}
	runnable := Use(prereq, options...)

	if cmd.Run != nil {
		Fail(fmt.Errorf("Run cannot be set on command"))
	}

	if body != nil {
		normalizedBody, err := castBodyFunction(body)
		if err != nil {
			Fail(err)
		}

		cmd.Run = func(_ *cobra.Command, args []string) {
			runnable.Run(func() error {
				return normalizedBody(args)
			})
		}
	}

	if parent != nil {
		parent.AddCommand(&cmd)
	}

	return &cmd
}
