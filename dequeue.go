package fuzzhelper

type dequeue[T any] struct {
	values []T
}

func newDequeue[T any]() *dequeue[T] {
	return &dequeue[T]{
		values: []T{},
	}
}

func (d *dequeue[T]) add(newValue T) {
	d.values = append(d.values, newValue)
}

func (d *dequeue[T]) addMany(newValues []T) {
	d.values = append(d.values, newValues...)
}

func (d *dequeue[T]) popFront() T {
	value := d.values[0]
	d.values = d.values[1:]
	return value
}

func (d *dequeue[T]) popBack() T {
	value := d.values[len(d.values)-1]
	d.values = d.values[:len(d.values)-1]
	return value
}

func (d *dequeue[T]) len() int {
	return len(d.values)
}
