package mm

type Vector[T any] struct {
	data []T
	cap  int
	len  int
}

func createVector[T any](len int, cap int) *Vector[T] {
	vector := Alloc[Vector[T]]()
	vector.cap = cap
	vector.len = len
	vector.data = AllocMany[T](vector.cap)

	return vector
}

// NewVector creates a new empty vector, if args not provided
// it will create an empty vector, if only one arg is provided
// it will init a vector with len and cap equal to the provided arg,
// if two args are provided it will init a vector with len = args[0] cap = args[1]
func NewVector[T any](args ...int) *Vector[T] {
	switch len(args) {
	case 0:
		return createVector[T](0, 1)
	case 1:
		return createVector[T](args[0], args[0])
	default:
		return createVector[T](args[0], args[1])
	}
}

// InitVector initializes a new vector with the T elements provided and sets
// it's len and cap to len(values)
func InitVector[T any](values ...T) *Vector[T] {
	vector := createVector[T](len(values), len(values))
	for _, v := range values {
		vector.data[0] = v
	}
	return vector
}

// Push pushes T to the vector, grows if needed.
func (v *Vector[T]) Push(value T) {
	if v.len == v.cap {
		v.data = Reallocate(v.data, v.cap*2)
		v.cap *= 2
	}

	v.data[v.len] = value
	v.len++
}

// Pop pops T from the vector
func (v *Vector[T]) Pop() T {
	v.len--
	return v.data[v.len]
}

// Len gets vector len
func (v *Vector[T]) Len() int {
	return v.len
}

// Cap gets vector capacity (underling memory length).
func (v *Vector[T]) Cap() int {
	return v.cap
}

// Slice gets a slice representing the vector
// CAUTION: don't append to this slice, this is only used
// if you want to loop on the vec elements
func (v *Vector[T]) Slice() []T {
	return v.data[:v.len]
}

// Deallocate deallocats the vector
func (v *Vector[T]) Free() {
	FreeMany(v.data)
	Free(v)
}
