package master

// Iterable is a helper to the concurrent execution tools RunTasksAndAccumulate
// and RunTasksAndAccumulateErrors. It provides a generic iterator that can be
// used to iterate over slices and maps of whatever types.
type Iterable[Idx comparable, Val any] interface {
	// Next increments the internal cursor to refer to the next object. It
	// returns true if another object exists or false if the end of iteration
	// has been reached.
	Next() bool

	// Id returns the ID value of the current value.
	Id() Idx

	// Value returns the value of the current value.
	Value() Val

	// Len reports the lengths of the iterated object.
	Len() int
}

// MapIterator is a generic implementation of Iterable for maps.
type MapIterator[Idx comparable, Val any] struct {
	is   map[Idx]Val
	keys []Idx
	idx  int
}

// NewMapIterator will create an iterator that iterates over a given map.
func NewMapIterator[Idx comparable, Val any](
	is map[Idx]Val,
) *MapIterator[Idx, Val] {
	keys := make([]Idx, 0, len(is))
	for k := range is {
		keys = append(keys, k)
	}

	return &MapIterator[Idx, Val]{is, keys, -1}
}

// Next returns false if there is no more work to do with this iterator. It
// returns true and increments the cursor pointer if there is more work to do.
func (i *MapIterator[Idx, Val]) Next() bool {
	if i.idx < len(i.keys) {
		i.idx++
	}
	return i.idx < len(i.keys)
}

// Value returns the current value of the key/value pair iteration.
func (i *MapIterator[Idx, Val]) Value() Val {
	return i.is[i.keys[i.idx]]
}

// Id returns the currnet key of the key/value pair iteration.
func (i *MapIterator[Idx, Val]) Id() Idx {
	return i.keys[i.idx]
}

// Len returns the len of the underlying map.
func (i *MapIterator[Idx, Val]) Len() int {
	return len(i.keys)
}

// SliceIterator provide a generic implementation of Iterable over slice
// objects. The Idx type is always int, in this case.
type SliceIterator[Val any] struct {
	is  []Val
	idx int
}

// NewSliceIterator creates a new iterator for the given slice.
func NewSliceIterator[Val any](
	is []Val,
) *SliceIterator[Val] {
	return &SliceIterator[Val]{is, -1}
}

// Next returns false if there are no more elements in the slice to process. It
// returns true and increments the index to operate upon if there is another
// element to process.
func (i *SliceIterator[Val]) Next() bool {
	if i.idx < len(i.is) {
		i.idx++
	}
	return i.idx < len(i.is)
}

// Value returns the value of the current slice element during iteration.
func (i *SliceIterator[Val]) Value() Val {
	return i.is[i.idx]
}

// Id returns the index of the current slice element during iteration.
func (i *SliceIterator[Val]) Id() int {
	return i.idx
}

// Len returns the len of the underlying slice.
func (i *SliceIterator[Val]) Len() int {
	return len(i.is)
}
