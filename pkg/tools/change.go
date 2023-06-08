package tools

func Point[T any](in T) *T {
	return &in
}
