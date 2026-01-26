package index

import "github.com/miguelm-revel/revelTools/collections"

type bkTree[T comparable] struct {
	value    string
	children map[int]*bkTree[T]
	content  T
}

func (b *bkTree[T]) search(value string, t int, nodes collections.Set[T]) {
	d0 := score(b.value, value)
	if d0 <= t {
		nodes.Add(b.content)
	}
	for k := range b.children {
		if k >= d0-t && k <= d0+t {
			b.children[k].search(value, t, nodes)
		}
	}
}

type BKTree[T comparable] struct {
	root *bkTree[T]
}

func (b *BKTree[T]) Insert(value string, content T) {
	node := &bkTree[T]{
		value:    value,
		children: make(map[int]*bkTree[T]),
		content:  content,
	}
	if b.root == nil {
		b.root = node
		return
	}
	n := b.root
	for {
		k := score(n.value, value)
		if child, ok := n.children[k]; !ok {
			n.children[k] = node
			return
		} else {
			n = child
		}
	}
}

func (b *BKTree[T]) Search(value string, t int) collections.Set[T] {
	ids := make(collections.Set[T])
	b.root.search(value, t, ids)
	return ids
}
