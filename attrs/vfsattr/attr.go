package vfsattr

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

const (
	ATTR_KEY   = "github.com/mandelsoft/vfs"
	ATTR_SHORT = "vfs"
)

func init() {
	ctxmgmt.RegisterAttributeType(ATTR_KEY, AttributeType{}, ATTR_SHORT)
}

type AttributeType struct{}

func (a AttributeType) Name() string {
	return ATTR_KEY
}

func (a AttributeType) Description() string {
	return `
*intern* (not via command line)
Virtual filesystem to use for command line context.
`
}

func (a AttributeType) Encode(v interface{}, marshaller runtime.Marshaler) ([]byte, error) {
	if _, ok := v.(vfs.FileSystem); !ok {
		return nil, fmt.Errorf("vfs.CachingFileSystem required")
	}
	return nil, nil
}

func (a AttributeType) Decode(data []byte, unmarshaller runtime.Unmarshaler) (interface{}, error) {
	return nil, errors.ErrNotSupported("decode attribute", ATTR_KEY)
}

////////////////////////////////////////////////////////////////////////////////

var _osfs = osfs.New()

func Get(ctx ctxmgmt.Context) vfs.FileSystem {
	v := ctx.GetAttributes().GetAttribute(ATTR_KEY)
	if v == nil {
		return _osfs
	}
	fs, _ := v.(vfs.FileSystem)
	return fs
}

func Set(ctx ctxmgmt.Context, fs vfs.FileSystem) {
	ctx.GetAttributes().SetAttribute(ATTR_KEY, fs)
}
