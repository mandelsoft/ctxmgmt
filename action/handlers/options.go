package handlers

import (
	"slices"

	"github.com/mandelsoft/ctxmgmt/action/api"
)

type (
	Option  = api.Option
	Options = api.Options
)

func NewOptions(opts ...Option) *Options {
	return api.NewOptions(opts...)
}

////////////////////////////////////////////////////////////////////////////////

type kind struct {
	action string
}

func ForAction(a string) Option {
	return kind{a}
}

func (o kind) ApplyActionHandlerOptionTo(opts *Options) {
	opts.Action = o.action
}

////////////////////////////////////////////////////////////////////////////////

type prio struct {
	prio int
}

func WithPrio(p int) Option {
	return prio{p}
}

func (o prio) ApplyActionHandlerOptionTo(opts *Options) {
	opts.Priority = o.prio
}

////////////////////////////////////////////////////////////////////////////////

type selectors struct {
	selectors []api.Selector
}

func ForSelectors(s ...api.Selector) Option {
	return selectors{s}
}

func (o selectors) ApplyActionHandlerOptionTo(opts *Options) {
	opts.Selectors = append(opts.Selectors, o.selectors...)
}

////////////////////////////////////////////////////////////////////////////////

type versions struct {
	versions []string
}

func WithVersions(vers ...string) Option {
	return versions{slices.Clone(vers)}
}

func (o versions) ApplyActionHandlerOptionTo(opts *Options) {
	opts.Versions = append(opts.Versions, o.versions...)
}
