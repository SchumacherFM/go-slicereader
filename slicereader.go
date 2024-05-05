package slicereader

import (
	"errors"
)

// EOS is the error returned by Read when no more element is available.
// Functions should return EOS only to signal a graceful end of input.
var EOS = errors.New("EOS")

type predicate[T any] func(v T) bool

// SliceReader supports reading a slice similar to an io.Reader.
type SliceReader[T any] struct {
	s []T
	i int64
}

// NewSliceReader returns a new SliceReader.
func NewSliceReader[T any](slice []T) *SliceReader[T] {
	return &SliceReader[T]{
		s: slice,
		i: 0,
	}
}

// Len returns the number of the unread elements of the slice.
func (sr *SliceReader[T]) Len() int {
	return int(sr.Size() - sr.i)
}

// Size returns the original length of the underlying slice.
// The returned value is always the same and is not affected by calls
// to any other method.
func (sr *SliceReader[T]) Size() int64 {
	return int64(len(sr.s))
}

// Read reads a single element of the slice or the EOS error will be returned.
func (sr *SliceReader[T]) Read() (e T, err error) {
	if sr.i >= sr.Size() {
		return e, EOS
	}
	e = sr.s[sr.i]
	sr.i++
	return
}

// ReadWhile reads the slice till the element before the given predicate
// function returns false. If the end of the slice is reached the EOS error will
// be returned.
func (sr *SliceReader[T]) ReadWhile(p predicate[T]) (s []T, err error) {
	for sr.i < sr.Size() {
		if !p(sr.s[sr.i]) {
			return s, nil
		}
		s = append(s, sr.s[sr.i])
		sr.i++
	}
	return s, EOS
}

// ReadUntil reads the slice till the element before the given predicate
// function returns true. If the end of the slice is reached the EOS error will
// be returned.
func (sr *SliceReader[T]) ReadUntil(p predicate[T]) (s []T, err error) {
	for sr.i < sr.Size() {
		if p(sr.s[sr.i]) {
			return s, nil
		}
		s = append(s, sr.s[sr.i])
		sr.i++
	}
	return s, EOS
}

// ReadWhileIncl reads the slice till including the element, the given predicate
// function returns false. If the end of the slice is reached the EOS error will
// be returned.
func (sr *SliceReader[T]) ReadWhileIncl(p predicate[T]) (s []T, err error) {
	for sr.i < sr.Size() {
		if !p(sr.s[sr.i]) {
			s = append(s, sr.s[sr.i])
			return s, nil
		}
		s = append(s, sr.s[sr.i])
		sr.i++
	}
	return s, EOS
}

// ReadUntilIncl reads the slice till including  the element, before the given
// predicate function returns true. If the end of the slice is reached the EOS
// error will be returned.
func (sr *SliceReader[T]) ReadUntilIncl(p predicate[T]) (s []T, err error) {
	for sr.i < sr.Size() {
		if p(sr.s[sr.i]) {
			s = append(s, sr.s[sr.i])
			return s, nil
		}
		s = append(s, sr.s[sr.i])
		sr.i++
	}
	return s, EOS
}
