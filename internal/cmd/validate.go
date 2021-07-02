package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type ValidateArgFn func(clause *kingpin.ArgClause, name string, predicateFn func(model *kingpin.ArgModel) (err error)) (err error)
type ValidateFlagFn func(clause *kingpin.ArgClause, name string, predicateFn func(model *kingpin.ArgModel) (err error)) (err error)
type ArgPredicateFn func(model *kingpin.ArgModel, value kingpin.Value) (err error)
type FlagPredicateFn func(model *kingpin.FlagModel, value kingpin.Value) (err error)

type Validator struct {
	argFns  map[string]ArgPredicateFn
	flagFns map[string]FlagPredicateFn
}

func NewValidator() *Validator {
	return &Validator{
		flagFns: map[string]FlagPredicateFn{},
		argFns:  map[string]ArgPredicateFn{},
	}
}

func (v *Validator) AddArgPredicateFn(name string, fn ArgPredicateFn) *Validator {
	v.argFns[name] = fn
	return v
}

func (v *Validator) AddFlagPredicateFn(name string, fn FlagPredicateFn) *Validator {
	v.flagFns[name] = fn
	return v
}

func (v *Validator) Validate(clause *kingpin.CmdClause) error {
	for _, arg := range clause.Model().Args {
		if argPredicateFn, ok := v.argFns[arg.Name]; ok {
			return argPredicateFn(arg, arg.Value)
		}
	}
	return nil
}
