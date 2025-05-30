package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gowebpki/jcs"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/set"
	"github.com/opencontainers/go-digest"
)

// Properties describes a set of name/value pairs.
type Properties map[string]string

// UniqueKey returns a unique key for the property set
func (p Properties) UniqueKey() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to produce digest: %w", err)
	}
	data, err = jcs.Transform(data)
	if err != nil {
		return nil, fmt.Errorf("failed to produce digest: %w", err)
	}
	return data, nil
}

// Digest returns the object digest of a Property set.
func (p Properties) Digest(algo ...digest.Algorithm) (*digest.Digest, error) {
	data, err := p.UniqueKey()
	if err != nil {
		return nil, err
	}
	d := general.OptionalDefaulted(digest.SHA256, algo...).FromBytes(data)
	return &d, nil
}

func (p Properties) SetNonEmptyValue(name, value string) {
	if value != "" {
		p[name] = value
	}
}

// Equals compares two identities.
func (p Properties) Equals(o Properties) bool {
	if len(p) != len(o) {
		return false
	}

	for k, v := range p {
		if v2, ok := o[k]; !ok || v != v2 {
			return false
		}
	}
	return true
}

// Match implements the selector interface.
func (p Properties) Match(obj map[string]string) (bool, error) {
	for k, v := range p {
		if obj[k] != v {
			return false, nil
		}
	}
	return true, nil
}

// Names returns the set of property names.
func (c Properties) Names() set.Set[string] {
	return set.KeySet(c)
}

// String returns a string representation.
func (c Properties) String() string {
	if c == nil {
		return "<none>"
	}
	//nolint: errchkjson // just a string map
	d, _ := json.Marshal(c)
	return string(d)
}

// Copy copies identity.
func (p Properties) Copy() Properties {
	if p == nil {
		return nil
	}
	n := Properties{}
	for k, v := range p {
		n[k] = v
	}
	return n
}
