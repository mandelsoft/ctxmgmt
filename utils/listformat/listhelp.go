package listformat

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/mandelsoft/goutils/maputils"
	"github.com/mandelsoft/goutils/stringutils"
)

type StringElementDescriptionList []string

func (l StringElementDescriptionList) Len() int                 { return len(l) / 2 }
func (l StringElementDescriptionList) Key(i int) string         { return l[2*i] }
func (l StringElementDescriptionList) Description(i int) string { return l[2*i+1] }

type StringElementList []string

func (l StringElementList) Len() int                 { return len(l) }
func (l StringElementList) Key(i int) string         { return l[i] }
func (l StringElementList) Description(i int) string { return "" }

func FormatList(def string, elems ...string) string {
	return FormatListElements(def, StringElementList(elems))
}

type maplist[K ~string, E any] struct {
	desc func(E) string
	keys []K
	m    map[K]E
}

func (l *maplist[K, E]) Len() int                 { return len(l.keys) }
func (l *maplist[K, E]) Key(i int) string         { return string(l.keys[i]) }
func (l *maplist[K, E]) Description(i int) string { return l.desc(l.m[l.keys[i]]) }

func FormatMapElements[K ~string, E any](def string, m map[K]E, desc ...func(E) string) string {
	if len(desc) == 0 || desc[0] == nil {
		desc = []func(E) string{StringDescription[E]}
	}
	keys := maputils.OrderedKeys(m)
	return FormatListElements(def, &maplist[K, E]{
		desc: desc[0],
		keys: keys,
		m:    m,
	})
}

type DescriptionSource interface {
	GetDescription() string
}

type DirectDescriptionSource interface {
	Description() string
}

func StringDescription[E any](e E) string {
	if d, ok := any(e).(DescriptionSource); ok {
		return d.GetDescription()
	}
	if d, ok := any(e).(DirectDescriptionSource); ok {
		return d.Description()
	}
	return fmt.Sprintf("%s", any(e))
}

type ListElements interface {
	Len() int
	Key(i int) string
	Description(i int) string
}

func FormatListElements(def string, elems ListElements) string {
	names := ""
	size := elems.Len()

	for i := 0; i < size; i++ {
		key := elems.Key(i)
		names = fmt.Sprintf("%s  - <code>%s</code>", names, key)
		if key == def {
			names += " (default)"
		}
		desc := elems.Description(i)
		if desc != "" {
			names += ": " + stringutils.IndentLines(desc, "    ", true)
			if strings.Contains(desc, "\n") {
				names += "\n"
			}
		}
		names += "\n"
	}
	return names
}

func FormatDescriptionList(def string, elems ...string) string {
	list := slices.Clone(elems)
	sort.Strings(list)

	out := ""
	for _, l := range list {
		if l != "" {
			out += "  - " + stringutils.IndentLines(l, "    ", true)
			if strings.Contains(l, "\n") {
				out += "\n"
			}
		}
		out += "\n"
	}
	return out
}
