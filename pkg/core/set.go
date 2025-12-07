package core

// Set is a generic set implementation
type Set[T comparable] map[T]struct{}
