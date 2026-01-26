package collections

import (
	"iter"
	"slices"
)

func Zip[X, Y any](a iter.Seq[X], b iter.Seq[Y]) iter.Seq2[X, Y] {
	return func(yield func(X, Y) bool) {
		nextA, stopA := iter.Pull(a)
		defer stopA()

		nextB, stopB := iter.Pull(b)
		defer stopB()

		for {
			ax, okA := nextA()
			by, okB := nextB()
			if !okA || !okB {
				return
			}
			if !yield(ax, by) {
				return
			}
		}
	}
}

func ZipSlice[X, Y any](a []X, b []Y) iter.Seq2[X, Y] {
	return Zip(slices.Values(a), slices.Values(b))
}
