package index

import (
	"github.com/miguelm-revel/revelTools/collections"
	"math"
	"math/rand/v2"
	"slices"
)

type hnsw[T comparable] struct {
	l       int
	vec     Vector
	content T
	nb      []collections.Set[*hnsw[T]]
}

type HNSW[T comparable] struct {
	ep *hnsw[T]
}

func (g *HNSW[T]) nearest(src collections.Set[*hnsw[T]], dst *hnsw[T]) *hnsw[T] {
	var nearest *hnsw[T]
	minDst := math.MaxFloat64
	for n := range src.Iter() {
		d := n.vec.Dist(dst.vec)
		if d < minDst {
			nearest = n
			minDst = d
		}
	}
	return nearest
}

func (g *HNSW[T]) furthest(src collections.Set[*hnsw[T]], dst *hnsw[T]) *hnsw[T] {
	var furthest *hnsw[T]
	maxDst := -math.MaxFloat64
	for n := range src.Iter() {
		d := n.vec.Dist(dst.vec)
		if d > maxDst {
			furthest = n
			maxDst = d
		}
	}
	return furthest
}

func (g *HNSW[T]) Insert(vec Vector, content any, M, MMax, EfConst int, ml float64) {
	if g.ep == nil {
		q := &hnsw[T]{
			l:       0,
			vec:     vec,
			content: content,
			nb:      make([]collections.Set[*hnsw[T]], 1),
		}
		q.nb[0] = make(collections.Set[*hnsw[T]])
		g.ep = q
	}
	var W collections.Set[*hnsw[T]]
	l := int(math.Floor(-math.Log(rand.Float64()) * ml))
	ep := []*hnsw[T]{g.ep}
	L := g.ep.l
	q := &hnsw[T]{
		l:       l,
		vec:     vec,
		content: content,
		nb:      make([]collections.Set[*hnsw[T]], l+1),
	}
	for i := 0; i <= l; i++ {
		q.nb[i] = make(collections.Set[*hnsw[T]])
	}
	for lc := g.ep.l; lc > l; lc-- {
		W = g.searchLayer(q, ep, 1, lc)
		ep = []*hnsw[T]{g.nearest(W, q)}
	}
	for lc := min(L, l); lc >= 0; lc-- {
		W = g.searchLayer(q, ep, EfConst, lc)
		neighbors := g.selectNeighbors(q, W, M)
		for _, n := range neighbors {
			n.nb[lc].Add(q)
			q.nb[lc].Add(n)
		}
		for _, e := range neighbors {
			eConn := e.nb[lc]
			if len(eConn) > MMax {
				eNewConn := g.selectNeighbors(e, eConn, MMax)
				e.nb[lc] = collections.NewSet(eNewConn)
			}
		}
		ep = make([]*hnsw[T], len(W))
		for i, n := range W.Iter2() {
			ep[i] = n
		}
	}
	if l > L {
		g.ep = q
	}
}

func (g *HNSW[T]) searchLayer(q *hnsw[T], ep []*hnsw[T], ef int, lc int) collections.Set[*hnsw[T]] {
	v := collections.NewSet(ep)
	C := collections.NewSet(ep)
	W := collections.NewSet(ep)
	for len(C) > 0 {
		c := g.nearest(C, q)
		C.Del(c)
		f := g.furthest(W, q)
		if c.vec.Dist(q.vec) > f.vec.Dist(q.vec) {
			break
		}
		for e := range c.nb[lc].Iter() {
			if !v.Has(e) {
				v.Add(e)
				f = g.furthest(W, q)
				if e.vec.Dist(q.vec) < f.vec.Dist(q.vec) {
					C.Add(e)
					W.Add(e)
					if len(W) > ef {
						W.Del(g.furthest(W, q))
					}
				}
			}
		}
	}
	return W
}

func (g *HNSW[T]) selectNeighbors(q *hnsw[T], C collections.Set[*hnsw[T]], M int) []*hnsw[T] {
	neighbors := make([]*hnsw[T], len(C))
	for i, c := range C.Iter2() {
		neighbors[i] = c
	}
	slices.SortFunc(neighbors, func(a, b *hnsw[T]) int {
		d1 := a.vec.Dist(q.vec)
		d2 := b.vec.Dist(q.vec)
		if d1 < d2 {
			return -1
		}
		if d1 > d2 {
			return 1
		}
		return 0
	})
	return neighbors[:min(len(C), M)]
}

func (g *HNSW[T]) KNNSearch(q *hnsw[T], K, ef int) []any {
	var W collections.Set[*hnsw[T]]
	ep := []*hnsw[T]{g.ep}
	L := g.ep.l
	for lc := L; lc >= 1; lc-- {
		W = g.searchLayer(q, ep, 1, lc)
		ep = []*hnsw[T]{g.nearest(W, q)}
	}
	W = g.searchLayer(q, ep, ef, 0)
	result := make([]any, 0)
	for _, doc := range g.selectNeighbors(q, W, K) {
		result = append(result, doc.content)
	}
	return result
}
