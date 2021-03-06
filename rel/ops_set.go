package rel

import (
	"sort"

	"github.com/arr-ai/frozen"
)

// Intersect returns every Value from a that is also in b.
func Intersect(a, b Set) Set {
	if ga, ok := a.(*genericSet); ok {
		if gb, ok := b.(*genericSet); ok {
			return &genericSet{set: ga.set.Intersection(gb.set)}
		}
	}
	return a.Where(func(v Value) bool { return b.Has(v) })
}

func NIntersect(a Set, bs ...Set) Set {
	for _, b := range bs {
		a = Intersect(a, b)
	}
	return a
}

// Union returns every Values that is in either input Set or both.
func Union(a, b Set) Set {
	if ga, ok := a.(*genericSet); ok {
		if gb, ok := b.(*genericSet); ok {
			return &genericSet{set: ga.set.Union(gb.set)}
		}
	}
	for e := b.Enumerator(); e.MoveNext(); {
		a = a.With(e.Current())
	}
	return a
}

func NUnion(sets ...Set) Set {
	result := None
	for _, s := range sets {
		result = Union(result, s)
	}
	return result
}

// Difference returns every Value from the first Set that is not in the second.
func Difference(a, b Set) Set {
	if ga, ok := a.(*genericSet); ok {
		if gb, ok := b.(*genericSet); ok {
			return &genericSet{set: ga.set.Difference(gb.set)}
		}
	}
	return a.Where(func(v Value) bool { return !b.Has(v) })
}

// SymmetricDifference returns Values in either Set, but not in both.
func SymmetricDifference(a, b Set) Set {
	if ga, ok := a.(*genericSet); ok {
		if gb, ok := b.(*genericSet); ok {
			return &genericSet{set: ga.set.SymmetricDifference(gb.set)}
		}
	}
	return Union(Difference(a, b), Difference(b, a))
}

// OrderBy returns a slice with the sets Values sorted by the given key.
func OrderBy(s Set, key func(v Value) (Value, error), less func(a, b Value) bool) ([]Value, error) {
	o := newOrderer(s.Count(), less)
	for i, e := 0, s.Enumerator(); e.MoveNext(); i++ {
		value := e.Current()
		o.values[i] = value
		var err error
		o.keys[i], err = key(value)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(o)
	return o.values, nil
}

type orderer struct {
	values []Value
	keys   []Value
	less   func(a, b Value) bool
}

func newOrderer(n int, less func(a, b Value) bool) *orderer {
	buf := make([]Value, 2*n)
	return &orderer{values: buf[:n], keys: buf[n:], less: less}
}

// Len is the number of elements in the collection.
func (o *orderer) Len() int {
	return len(o.values)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (o *orderer) Less(i, j int) bool {
	return o.less(o.keys[i], o.keys[j])
}

// Swap swaps the elements with indexes i and j.
func (o *orderer) Swap(i, j int) {
	o.values[i], o.values[j] = o.values[j], o.values[i]
	o.keys[i], o.keys[j] = o.keys[j], o.keys[i]
}

// PowerSet computes the power set of a set.
func PowerSet(s Set) Set {
	if gs, ok := s.(*genericSet); ok {
		var sb frozen.SetBuilder
		for i := gs.set.Powerset().Range(); i.Next(); {
			sb.Add(&genericSet{set: i.Value().(frozen.Set)})
		}
		return &genericSet{set: sb.Finish()}
	}
	result := NewSet(None)
	for e := s.Enumerator(); e.MoveNext(); {
		c := e.Current()
		newSets := NewSet()
		for s := result.Enumerator(); s.MoveNext(); {
			newSets = newSets.With(s.Current().(Set).With(c))
		}
		result = Union(result, newSets)
	}
	return result
}
