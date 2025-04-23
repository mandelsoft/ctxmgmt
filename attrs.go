package datacontext

import (
	"io"
	"sort"
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/finalizer"
	"github.com/mandelsoft/goutils/general"

	"github.com/mandelsoft/datacontext/utils"
	"github.com/mandelsoft/datacontext/utils/runtime"
)

type AttributeType interface {
	Name() string
	Decode(data []byte, unmarshaler runtime.Unmarshaler) (interface{}, error)
	Encode(v interface{}, marshaller runtime.Marshaler) ([]byte, error)
	Description() string
}

// Converter is an optional interface an AttributeType can implement to
// normalize an attribute value. It is called by the Attributes.SetAttribute
// method.
type Converter interface {
	Convert(interface{}) (interface{}, error)
}

type AttributeScheme interface {
	Register(name string, typ AttributeType, short ...string) error

	Decode(attr string, data []byte, unmarshaler runtime.Unmarshaler) (interface{}, error)
	Encode(attr string, v interface{}, marshaller runtime.Marshaler) ([]byte, error)
	Convert(attr string, v interface{}) (interface{}, error)
	GetType(attr string) (AttributeType, error)

	AddKnownTypes(scheme AttributeScheme)
	Shortcuts() utils.Properties
	KnownTypes() KnownTypes
	KnownTypeNames() []string
}

var DefaultAttributeScheme = NewDefaultAttributeScheme()

// KnownTypes is a set of known type names mapped to appropriate object decoders.
type KnownTypes map[string]AttributeType

// Copy provides a copy of the actually known types.
func (t KnownTypes) Copy() KnownTypes {
	n := KnownTypes{}
	for k, v := range t {
		n[k] = v
	}
	return n
}

// TypeNames return a sorted list of known type names.
func (t KnownTypes) TypeNames() []string {
	types := make([]string, 0, len(t))
	for t := range t {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

type defaultScheme struct {
	lock  sync.RWMutex
	types KnownTypes
	short utils.Properties
}

func NewDefaultAttributeScheme() AttributeScheme {
	return &defaultScheme{
		types: KnownTypes{},
		short: utils.Properties{},
	}
}

func (d *defaultScheme) AddKnownTypes(s AttributeScheme) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for k, v := range s.KnownTypes() {
		d.types[k] = v
	}
	for k, v := range s.Shortcuts() {
		d.short[k] = v
	}
}

func (d *defaultScheme) KnownTypes() KnownTypes {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.types.Copy()
}

func (d *defaultScheme) Shortcuts() utils.Properties {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.short.Copy()
}

// KnownTypeNames return a sorted list of known type names.
func (d *defaultScheme) KnownTypeNames() []string {
	d.lock.RLock()
	defer d.lock.RUnlock()
	types := make([]string, 0, len(d.types))
	for t := range d.types {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

func RegisterAttributeType(name string, typ AttributeType, short ...string) error {
	return DefaultAttributeScheme.Register(name, typ, short...)
}

func (d *defaultScheme) Register(name string, typ AttributeType, short ...string) error {
	if typ == nil {
		return errors.Newf("type object must be given")
	}
	if name == "" {
		return errors.Newf("name must be given")
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	d.types[name] = typ
	for _, s := range short {
		d.short[s] = name
	}
	return nil
}

func (d *defaultScheme) getType(attr string) AttributeType {
	if s, ok := d.short[attr]; ok {
		attr = s
	}
	return d.types[attr]
}

func (d *defaultScheme) GetType(attr string) (AttributeType, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	t := d.getType(attr)
	if t == nil {
		return nil, errors.ErrUnknown("attribute", attr)
	}
	return t, nil
}

func (d *defaultScheme) Convert(attr string, value interface{}) (interface{}, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	t := d.getType(attr)
	if t == nil {
		return value, errors.ErrUnknown("attribute", attr)
	}
	if c, ok := t.(Converter); ok {
		return c.Convert(value)
	}
	return value, nil
}

func (d *defaultScheme) Encode(attr string, value interface{}, marshaler runtime.Marshaler) ([]byte, error) {
	if marshaler == nil {
		marshaler = runtime.DefaultJSONEncoding
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	t := d.getType(attr)
	if t == nil {
		return nil, errors.ErrUnknown("attribute", attr)
	}
	return t.Encode(value, marshaler)
}

func (d *defaultScheme) Decode(attr string, data []byte, unmarshaler runtime.Unmarshaler) (interface{}, error) {
	if unmarshaler == nil {
		unmarshaler = runtime.DefaultJSONEncoding
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	t := d.getType(attr)
	if t == nil {
		return nil, errors.ErrUnknown("attribute", attr)
	}
	return t.Decode(data, unmarshaler)
}

type DefaultAttributeType struct{}

func (_ DefaultAttributeType) Encode(v interface{}, marshaller runtime.Marshaler) ([]byte, error) {
	return marshaller.Marshal(v)
}

////////////////////////////////////////////////////////////////////////////////

type _attributes struct {
	sync.RWMutex
	id         uint64
	ctx        Context
	parent     Attributes
	updater    *Updater
	attributes map[string]interface{}
}

var _ Attributes = &_attributes{}

func NewAttributes(ctx Context, parent Attributes, updater *Updater) Attributes {
	return newAttributes(ctx, parent, updater)
}

func newAttributes(ctx Context, parent Attributes, updater *Updater) *_attributes {
	return &_attributes{
		id:         attrsrange.NextId(),
		ctx:        ctx,
		parent:     parent,
		updater:    updater,
		attributes: map[string]interface{}{},
	}
}

func (c *_attributes) Finalize() error {
	list := errors.ErrListf("finalizing attributes")
	for n, a := range c.attributes {
		if f, ok := a.(finalizer.Finalizable); ok {
			list.Addf(nil, f.Finalize(), "attribute %s", n)
		}
	}
	return list.Result()
}

func (c *_attributes) GetAttribute(name string, def ...interface{}) interface{} {
	if *c.updater != nil {
		(*c.updater).Update()
	}
	c.RLock()
	defer c.RUnlock()
	if a := c.attributes[name]; a != nil {
		return a
	}
	if c.parent != nil {
		if a := c.parent.GetAttribute(name); a != nil {
			return a
		}
	}
	return general.Optional(def...)
}

func (c *_attributes) SetEncodedAttribute(name string, data []byte, unmarshaller runtime.Unmarshaler) error {
	s := DefaultAttributeScheme.Shortcuts()[name]
	if s != "" {
		name = s
	}
	v, err := DefaultAttributeScheme.Decode(name, data, unmarshaller)
	if err != nil {
		return err
	}
	c.SetAttribute(name, v)
	return nil
}

func (c *_attributes) setAttribute(name string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	_, err := DefaultAttributeScheme.Encode(name, value, nil)
	if err != nil && !errors.IsErrUnknownKind(err, "attribute") {
		return err
	}
	old := c.attributes[name]
	if old != nil && old != value {
		if c, ok := old.(io.Closer); ok {
			c.Close()
		}
	}
	value, err = DefaultAttributeScheme.Convert(name, value)
	if err != nil && !errors.IsErrUnknownKind(err, "attribute") {
		return err
	}
	c.attributes[name] = value
	return nil
}

func (c *_attributes) SetAttribute(name string, value interface{}) error {
	err := c.setAttribute(name, value)
	if err == nil {
		if *c.updater != nil {
			(*c.updater).Update()
		}
	}
	return err
}

func (c *_attributes) getOrCreateAttribute(name string, creator AttributeFactory) interface{} {
	c.Lock()
	defer c.Unlock()
	if v := c.attributes[name]; v != nil {
		return v
	}
	if c.parent != nil {
		if v := c.parent.GetAttribute(name); v != nil {
			return v
		}
	}
	v := creator(c.ctx)
	c.attributes[name] = v
	return v
}

func (c *_attributes) GetOrCreateAttribute(name string, creator AttributeFactory) interface{} {
	r := c.getOrCreateAttribute(name, creator)
	if *c.updater != nil {
		(*c.updater).Update()
	}
	return r
}
