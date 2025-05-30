package runtime

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/mandelsoft/ctxmgmt/utils/errkind"
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/modern-go/reflect2"
)

var (
	typeTypedObject = reflect.TypeOf((*TypedObject)(nil)).Elem()
	typeUnknown     = reflect.TypeOf((*Unknown)(nil)).Elem()
)

type (
	// TypedObjectDecoder is able to provide an effective typed object for some
	// serilaized form. The technical deserialization is done by an Unmarshaler.
	TypedObjectDecoder[T TypedObject] interface {
		Decode(data []byte, unmarshaler Unmarshaler) (T, error)
	}
	_TypedObjectDecoder[T TypedObject] interface {
		TypedObjectDecoder[T]
	}
)

// TypedObjectEncoder is able to provide a versioned representation of
// an effective TypedObject.
type TypedObjectEncoder[T TypedObject] interface {
	Encode(T, Marshaler) ([]byte, error)
}

type DirectDecoder[T TypedObject] struct {
	proto reflect.Type
}

var _ TypedObjectDecoder[TypedObject] = &DirectDecoder[TypedObject]{}

func MustNewDirectDecoder[T TypedObject](proto T) *DirectDecoder[T] {
	d, err := NewDirectDecoder[T](proto)
	if err != nil {
		panic(err)
	}
	return d
}

func NewDirectDecoder[T TypedObject](proto T) (*DirectDecoder[T], error) {
	t := MustProtoType(proto)
	if !reflect.PointerTo(t).Implements(typeTypedObject) {
		return nil, errors.Newf("object interface %T: must implement TypedObject", proto)
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.Newf("prototype %q must be a struct", t)
	}
	return &DirectDecoder[T]{
		proto: t,
	}, nil
}

func (d *DirectDecoder[T]) CreateInstance() T {
	return reflect.New(d.proto).Interface().(T)
}

func (d *DirectDecoder[T]) Decode(data []byte, unmarshaler Unmarshaler) (T, error) {
	var zero T
	inst := d.CreateInstance()
	err := unmarshaler.Unmarshal(data, inst)
	if err != nil {
		return zero, err
	}

	return inst, nil
}

func (d *DirectDecoder[T]) Encode(obj T, marshaler Marshaler) ([]byte, error) {
	return marshaler.Marshal(obj)
}

// KnownTypes is a set of known type names mapped to appropriate object decoders.
type KnownTypes[T TypedObject, R TypedObjectDecoder[T]] map[string]R

// Copy provides a copy of the actually known types.
func (t KnownTypes[T, R]) Copy() KnownTypes[T, R] {
	n := KnownTypes[T, R]{}
	for k, v := range t {
		n[k] = v
	}
	return n
}

// TypeNames return a sorted list of known type names.
func (t KnownTypes[T, R]) TypeNames() []string {
	types := make([]string, 0, len(t))
	for t := range t {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// Unknown is the interface to be implemented by
// representations on an unknown, but nevertheless decoded specification
// of a typed object.
type Unknown interface {
	IsUnknown() bool
}

func IsUnknown(o TypedObject) bool {
	if reflect2.IsNil(o) {
		return true
	}
	if u, ok := o.(Unknown); ok {
		return u.IsUnknown()
	}
	return false
}

type (
	// Scheme is the interface to describe a set of object types
	// that implement a dedicated interface.
	// As such it knows about the desired interface of the instances
	// and can validate it. Additionally, it provides an implementation
	// for generic unstructured objects that can be used to decode
	// any serialized from of object candidates and provide the
	// effective type.
	Scheme[T TypedObject, R TypedObjectDecoder[T]] interface {
		SchemeCommon
		KnownTypesProvider[T, R]
		TypedObjectEncoder[T]
		TypedObjectDecoder[T]

		BaseScheme() Scheme[T, R] // Go does not support an additional type parameter S Scheme[T,S] to return the correct type here

		AddKnownTypes(scheme KnownTypesProvider[T, R])
		RegisterByDecoder(typ string, decoder R) error

		ValidateInterface(object T) error
		CreateUnstructured() T
		Convert(object TypedObject) (T, error)
		GetDecoder(otype string) R
		EnforceDecode(data []byte, unmarshaler Unmarshaler) (T, error)
	}
	_Scheme[T TypedObject, R TypedObjectDecoder[T]] interface { // cannot omit nesting, because Goland does not accept it
		Scheme[T, R]
	}
)

type KnownTypesProvider[T TypedObject, R TypedObjectDecoder[T]] interface {
	KnownTypes() KnownTypes[T, R]
}

type SchemeCommon interface {
	KnownTypeNames() []string
}

type defaultScheme[T TypedObject, R TypedObjectDecoder[T]] struct {
	lock           sync.RWMutex
	base           Scheme[T, R]
	instance       reflect.Type
	unstructured   reflect.Type
	defaultdecoder TypedObjectDecoder[T]
	acceptUnknown  bool
	types          KnownTypes[T, R]
}

var _ Scheme[VersionedTypedObject, TypedObjectDecoder[VersionedTypedObject]] = (*defaultScheme[VersionedTypedObject, TypedObjectDecoder[VersionedTypedObject]])(nil)

func MustNewDefaultScheme[T TypedObject, R TypedObjectDecoder[T]](protoUnstr Unstructured, acceptUnknown bool, defaultdecoder TypedObjectDecoder[T], base ...Scheme[T, R]) Scheme[T, R] {
	return general.Must(NewDefaultScheme[T](protoUnstr, acceptUnknown, defaultdecoder, base...))
}

func NewScheme[T TypedObject, R TypedObjectDecoder[T]](base ...Scheme[T, R]) Scheme[T, R] {
	s, _ := NewDefaultScheme[T](nil, false, nil, base...)
	return s
}

func NewDefaultScheme[T TypedObject, R TypedObjectDecoder[T]](protoUnstr Unstructured, acceptUnknown bool, defaultdecoder TypedObjectDecoder[T], base ...Scheme[T, R]) (Scheme[T, R], error) {
	var err error

	var protoIfce T
	it := reflect.TypeOf(&protoIfce)
	for it.Kind() == reflect.Ptr {
		it = it.Elem()
	}

	var ut reflect.Type
	if acceptUnknown {
		ut, err = ProtoType(protoUnstr)
		if err != nil {
			return nil, errors.Wrapf(err, "unstructured prototype %T", protoUnstr)
		}
		if !reflect.PointerTo(ut).Implements(it) {
			return nil, fmt.Errorf("unstructured type %T must implement %T to be acceptale as unknown result", protoUnstr, &protoIfce)
		}
		if !reflect.PointerTo(ut).Implements(typeUnknown) {
			return nil, fmt.Errorf("unstructured type %T must implement Unknown to be acceptable as unknown result", protoUnstr)
		}
	}

	b := general.Optional(base...)
	return &defaultScheme[T, R]{
		base:          b,
		instance:      it,
		unstructured:  ut,
		types:         KnownTypes[T, R]{},
		acceptUnknown: acceptUnknown,
	}, nil
}

func (d *defaultScheme[T, R]) BaseScheme() Scheme[T, R] {
	return d.base
}

func (d *defaultScheme[T, R]) AddKnownTypes(s KnownTypesProvider[T, R]) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for k, v := range s.KnownTypes() {
		d.types[k] = v
	}
}

func (d *defaultScheme[T, R]) KnownTypes() KnownTypes[T, R] {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.base == nil {
		return d.types.Copy()
	}
	kt := d.base.KnownTypes()
	for n, t := range d.types {
		kt[n] = t
	}
	return kt
}

// KnownTypeNames return a sorted list of known type names.
func (d *defaultScheme[T, R]) KnownTypeNames() []string {
	d.lock.RLock()
	defer d.lock.RUnlock()

	types := make([]string, 0, len(d.types))
	for t := range d.types {
		types = append(types, t)
	}
	if d.base != nil {
		types = append(types, d.base.KnownTypeNames()...)
	}
	sort.Strings(types)
	return types
}

func (d *defaultScheme[T, R]) RegisterByDecoder(typ string, decoder R) error {
	if reflect2.IsNil(decoder) {
		return errors.Newf("decoder must be given")
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	d.types[typ] = decoder
	return nil
}

func (d *defaultScheme[T, R]) ValidateInterface(object T) error {
	t := reflect.TypeOf(object)
	if !t.Implements(d.instance) {
		return errors.Newf("object type %q does not implement required instance interface %q", t, d.instance)
	}
	return nil
}

func (d *defaultScheme[T, R]) GetDecoder(typ string) R {
	d.lock.RLock()
	defer d.lock.RUnlock()
	decoder := d.types[typ]
	if reflect2.IsNil(decoder) && d.base != nil {
		decoder = d.base.GetDecoder(typ)
	}
	return decoder
}

func (d *defaultScheme[T, R]) CreateUnstructured() T {
	var _nil T
	if d.unstructured == nil {
		return _nil
	}
	return reflect.New(d.unstructured).Interface().(T)
}

func (d *defaultScheme[T, R]) Encode(obj T, marshaler Marshaler) ([]byte, error) {
	if marshaler == nil {
		marshaler = DefaultYAMLEncoding
	}
	decoder := d.GetDecoder(obj.GetType())
	if encoder, ok := generics.TryCast[TypedObjectEncoder[T]](decoder); ok {
		return encoder.Encode(obj, marshaler)
	}
	return marshaler.Marshal(obj)
}

func (d *defaultScheme[T, R]) Decode(data []byte, unmarshal Unmarshaler) (T, error) {
	var _nil T

	if unmarshal == nil {
		unmarshal = DefaultYAMLEncoding
	}

	un := d.CreateUnstructured()
	t := ""
	if reflect2.IsNil(un) {
		t, _ = DefaultProvider{}.GetTypeFor(data, unmarshal)
	} else {
		err := unmarshal.Unmarshal(data, un)
		if err != nil {
			return _nil, errors.Wrapf(err, "cannot unmarshal unstructured")
		}
		t = un.GetType()
	}

	if t == "" {
		return _nil, errors.Newf("no type found")
	}
	decoder := d.GetDecoder(t)
	if reflect2.IsNil(decoder) {
		if d.defaultdecoder != nil {
			o, err := d.defaultdecoder.Decode(data, unmarshal)
			if err == nil {
				if !reflect2.IsNil(o) {
					return o, nil
				}
			} else if !errors.IsErrUnknownKind(err, errkind.KIND_OBJECTTYPE) {
				return _nil, err
			}
		}
		if d.acceptUnknown && !reflect2.IsNil(un) {
			return un, nil
		}
		return _nil, errors.ErrUnknown(errkind.KIND_OBJECTTYPE, t)
	}
	return decoder.Decode(data, unmarshal)
}

func (d *defaultScheme[T, R]) EnforceDecode(data []byte, unmarshal Unmarshaler) (T, error) {
	var _nil T

	un := d.CreateUnstructured()
	if unmarshal == nil {
		unmarshal = DefaultYAMLEncoding.Unmarshaler
	}
	err := unmarshal.Unmarshal(data, un)
	if err != nil {
		return _nil, errors.Wrapf(err, "cannot unmarshal unstructured")
	}
	if un.GetType() == "" {
		if d.acceptUnknown {
			return un, nil
		}
		return un, errors.Newf("no type found")
	}
	decoder := d.GetDecoder(un.GetType())
	if reflect2.IsNil(decoder) {
		if d.defaultdecoder != nil {
			o, err := d.defaultdecoder.Decode(data, unmarshal)
			if err == nil {
				return o, nil
			}
			if !errors.IsErrUnknownKind(err, errkind.KIND_OBJECTTYPE) {
				return un, err
			}
		}
		if d.acceptUnknown {
			return un, nil
		}
		return un, errors.ErrUnknown(errkind.KIND_OBJECTTYPE, un.GetType())
	}
	o, err := decoder.Decode(data, unmarshal)
	if err != nil {
		return un, err
	}
	return o, err
}

func (d *defaultScheme[T, R]) Convert(o TypedObject) (T, error) {
	var _nil T

	if o.GetType() == "" {
		return _nil, errors.Newf("no type found")
	}

	if u, ok := o.(T); ok {
		return u, nil
	}

	if u, ok := o.(Unstructured); ok {
		raw, err := u.GetRaw()
		if err != nil {
			return _nil, err
		}
		return d.Decode(raw, DefaultJSONEncoding)
	}

	data, err := json.Marshal(o)
	if err != nil {
		return _nil, err
	}
	decoder := d.GetDecoder(o.GetType())
	if reflect2.IsNil(decoder) {
		if d.defaultdecoder != nil {
			object, err := d.defaultdecoder.Decode(data, DefaultJSONEncoding)
			if err == nil {
				return object, nil
			}
			if !errors.IsErrUnknownKind(err, errkind.KIND_OBJECTTYPE) {
				return _nil, err
			}
		}
		return _nil, errors.ErrUnknown(errkind.KIND_OBJECTTYPE, o.GetType())
	}
	r, err := decoder.Decode(data, DefaultJSONEncoding)
	if err != nil {
		return _nil, err
	}
	if reflect.TypeOf(r) == reflect.TypeOf(o) {
		return o.(T), nil
	}
	return r, nil
}

////////////////////////////////////////////////////////////////////////////////

// TypeScheme is a scheme based on Types instead of decoders.
type TypeScheme[T TypedObject, R TypedObjectType[T]] interface {
	Scheme[T, R]

	Register(typ R)
	GetType(name string) R
}

type defaultTypeScheme[T TypedObject, R TypedObjectType[T]] struct {
	_Scheme[T, R]
}

func MustNewDefaultTypeScheme[T TypedObject, R TypedObjectType[T]](protoUnstr Unstructured, acceptUnknown bool, defaultdecoder TypedObjectDecoder[T], base ...TypeScheme[T, R]) TypeScheme[T, R] {
	return general.Must(NewDefaultTypeScheme[T, R](protoUnstr, acceptUnknown, defaultdecoder, base...))
}

func NewTypeScheme[T TypedObject, R TypedObjectType[T]](base ...TypeScheme[T, R]) TypeScheme[T, R] {
	s, _ := NewDefaultTypeScheme[T](nil, false, nil, base...)
	return s
}

func NewDefaultTypeScheme[T TypedObject, R TypedObjectType[T]](protoUnstr Unstructured, acceptUnknown bool, defaultdecoder TypedObjectDecoder[T], base ...TypeScheme[T, R]) (TypeScheme[T, R], error) {
	s, err := NewDefaultScheme[T](protoUnstr, acceptUnknown, defaultdecoder, generics.Cast[Scheme[T, R]](general.Optional(base...)))
	if err != nil {
		return nil, err
	}
	return &defaultTypeScheme[T, R]{s}, nil
}

func (s *defaultTypeScheme[T, R]) Register(t R) {
	s.RegisterByDecoder(t.GetType(), t)
}

func (s *defaultTypeScheme[T, R]) GetType(name string) R {
	return generics.Cast[R](s.GetDecoder(name))
}
