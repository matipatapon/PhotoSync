package helper

type List[T any] struct {
	items []T
}

func (l *List[T]) Append(val T) {
	l.items = append(l.items, val)
}

func (l *List[T]) PopFirst() T {
	val := l.items[0]
	l.items = l.items[1:len(l.items)]
	return val
}

func (l *List[T]) Length() int {
	return len(l.items)
}
