package datacontext

import (
	"context"

	"github.com/mandelsoft/datacontext/action/handlers"
	ctxlog "github.com/mandelsoft/datacontext/logging"
	"github.com/mandelsoft/datacontext/utils/refmgmt"
	"github.com/mandelsoft/datacontext/utils/runtime"
	"github.com/mandelsoft/datacontext/utils/runtimefinalizer"
	"github.com/mandelsoft/goutils/finalizer"
)

type ContextIdentity = runtimefinalizer.ObjectIdentity

type ContextProvider interface {
	// AttributesContext returns the shared attributes
	AttributesContext() AttributesContext
}

// Delegates is the interface for common
// Context features, which might be delegated
// to aggregated contexts.
type Delegates interface {
	ctxlog.LogProvider
	handlers.ActionsProvider
}

type ContextBinder interface {
	// BindTo binds the context to a context.Context and makes it
	// retrievable by a ForContext method
	BindTo(ctx context.Context) context.Context
}

// Context describes a common interface for a data context used for a dedicated
// purpose.
// Such has a type and always specific attribute store.
// Every Context can be bound to a context.Context.
type Context interface {
	ContextBinder
	ContextProvider
	Delegates

	IsIdenticalTo(Context) bool

	// GetType returns the context type
	GetType() string
	GetId() ContextIdentity

	GetAttributes() Attributes

	Finalize() error
	Finalizer() *finalizer.Finalizer
}

type InternalContext interface {
	Context
	runtimefinalizer.RecorderProvider
	GetKey() interface{}
	GetAllocatable() refmgmt.Allocatable
}

type Attributes interface {
	finalizer.Finalizable

	GetAttribute(name string, def ...interface{}) interface{}
	SetAttribute(name string, value interface{}) error
	SetEncodedAttribute(name string, data []byte, unmarshaller runtime.Unmarshaler) error
	GetOrCreateAttribute(name string, creator AttributeFactory) interface{}
}

// Updater is the interface for contexts and other objects
// supporting a configuration update.
type Updater interface {
	Update() error
}

type UpdateFunc func() error

func (u UpdateFunc) Update() error {
	return u()
}

// AttributeFactory is used to atomically create a new attribute for a context.
type AttributeFactory func(Context) interface{}

// AttributesContext is the interface of the root context.
type AttributesContext interface {
	Context

	IsAttributesContext() bool
	AttributesContext() AttributesContext

	BindTo(ctx context.Context) context.Context
}
