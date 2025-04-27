package internal

import (
	"sync"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/maputils"
	"github.com/mandelsoft/goutils/sliceutils"
)

const KIND_CONFIGAPPLIER = "config applier"

// ConfigApplier is a generic handler able to
// apply a configuration to a configuration target.
// It can be used by generic config objects
// to attach a concrete configuration behaviour
// to its generic config data.
type ConfigApplier interface {
	ApplyConfigTo(ctx Context, cfg, tgt interface{}) error
}

type ConfigApplierRegistry interface {
	Register(name string, applier ConfigApplier)
	AddKnown(o ConfigApplierRegistry)
	Get(name string) ConfigApplier
	Names() []string
}

var DefaultConfigApplierRegistry ConfigApplierRegistry = NewConfigApplierRegistry()

type _ConfigApplierRegistry struct {
	lock     sync.Mutex
	base     ConfigApplierRegistry
	appliers map[string]ConfigApplier
}

func NewConfigApplierRegistry(base ...ConfigApplierRegistry) ConfigApplierRegistry {
	return &_ConfigApplierRegistry{base: general.Optional(base...), appliers: make(map[string]ConfigApplier)}
}
func (r *_ConfigApplierRegistry) Register(name string, applier ConfigApplier) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.appliers[name] = applier
}

func (r *_ConfigApplierRegistry) AddKnown(o ConfigApplierRegistry) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, n := range o.Names() {
		r.Register(n, o.Get(n))
	}
}

func (r *_ConfigApplierRegistry) Names() []string {
	r.lock.Lock()
	defer r.lock.Unlock()
	names := maputils.OrderedKeys(r.appliers)
	if r.base != nil {
		names = sliceutils.AppendUnique(names, r.base.Names()...)
	}
	return names
}

func (r *_ConfigApplierRegistry) Get(name string) ConfigApplier {
	r.lock.Lock()
	defer r.lock.Unlock()
	applier := r.appliers[name]
	if applier == nil && r.base != nil {
		applier = r.base.Get(name)
	}
	return applier
}
