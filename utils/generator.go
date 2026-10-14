package utils

type Generator[T any] struct {
	data    []T
	pointer int
}

func NewGenerator[T any](data []T) *Generator[T] {
	return &Generator[T]{
		data:    data,
		pointer: 0,
	}
}

func (g *Generator[T]) next() (value T, isEnd bool) {
	value, isEnd = g.peek()
	if !isEnd {
		g.pointer++
	}
	return
}

func (g Generator[T]) peek() (T, bool) {
	length := len(g.data)
	if g.pointer >= length {
		return g.data[length-1], true
	}
	return g.data[g.pointer], false
}
