package tmpcache

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/attrs/vfsattr"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

const (
	ATTR_KEY   = "github.com/mandelsoft/tempblobcache"
	ATTR_SHORT = "blobcache"
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
*string* Folder name for temporary blob cache
The temporary blob cache is used to accessing large blobs from remote systems.
The are temporarily stored in the filesystem, instead of the memory, to avoid
blowing up the memory consumption.
`
}

func (a AttributeType) Encode(v interface{}, marshaller runtime.Marshaler) ([]byte, error) {
	if a, ok := v.(*Attribute); !ok {
		return nil, fmt.Errorf("temppcache attribute")
	} else {
		return []byte(a.Path), nil
	}
}

func (a AttributeType) Decode(data []byte, unmarshaller runtime.Unmarshaler) (interface{}, error) {
	var s string
	err := runtime.DefaultYAMLEncoding.Unmarshal(data, &s)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid attribute value for %s", ATTR_KEY)
	}
	return &Attribute{
		Path: s,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

type Attribute struct {
	Path       string
	Filesystem vfs.FileSystem
}

func New(path string, fss ...vfs.FileSystem) *Attribute {
	fs := general.OptionalDefaulted[vfs.FileSystem](osfs.OsFs, fss...)
	if path == "" {
		path = fs.FSTempDir()
	}
	return &Attribute{
		Path:       path,
		Filesystem: fs,
	}
}

func (a *Attribute) CreateTempFile(pat string) (vfs.File, error) {
	err := a.Filesystem.MkdirAll(a.Path, 0o777)
	if err != nil {
		return nil, err
	}
	return vfs.TempFile(a.Filesystem, a.Path, pat)
}

////////////////////////////////////////////////////////////////////////////////

func Get(ctx ctxmgmt.Context) *Attribute {
	var v interface{}
	var fs vfs.FileSystem

	if ctx != nil {
		v = ctx.GetAttributes().GetAttribute(ATTR_KEY)
		fs = vfsattr.Get(ctx)
	}

	if v != nil {
		a := v.(*Attribute)
		if a.Filesystem == nil {
			a.Filesystem = fs
		}
		return a
	}
	return New("", fs)
}

func Set(ctx ctxmgmt.Context, a *Attribute) {
	ctx.GetAttributes().SetAttribute(ATTR_KEY, a)
}
