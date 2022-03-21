package chewbyte

// Don't really look at this file. I wrote it mostly to practice generics in Go,
// most of the stuff here is not really needed

type Set[T comparable] struct {
	container map[T]bool
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		container: make(map[T]bool, 0),
	}
}

func (s *Set[T]) Has(item T) bool {
	_, ok := s.container[item]
	return ok
}

func (s *Set[T]) Insert(item T) {
	s.container[item] = true
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.container) == 0
}

func (s *Set[T]) GetSlice() []T {
	r := make([]T, 0)

	for k := range s.container {
		r = append(r, k)
	}

	return r
}

func (s *Set[T]) Pop() (T, bool) {
	for k := range s.container {
		delete(s.container, k)
		return k, true
	}

	var zero T
	return zero, false

}

type VisitableQueue[T comparable] struct {
	visited Set[T]
	pending Set[T]
}

func NewVisitableQueue[T comparable]() VisitableQueue[T] {
	return VisitableQueue[T]{
		visited: NewSet[T](),
		pending: NewSet[T](),
	}
}

func (q *VisitableQueue[T]) IsEmpty() bool {
	return q.pending.IsEmpty()
}

func (q *VisitableQueue[T]) Put(item T) {
	// Prevent from revisiting already visited items
	if q.visited.Has(item) {
		return
	}

	q.pending.Insert(item)
}

func (q *VisitableQueue[T]) Import(items []T) {
	for _, v := range items {
		q.Put(v)
	}
}

func (q *VisitableQueue[T]) Get() (T, bool) {
	element, ok := q.pending.Pop()
	if !ok {
		var zero T
		return zero, false
	}

	q.visited.Insert(element)
	return element, true
}
