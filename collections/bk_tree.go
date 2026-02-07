package collections

const (
	MATCH    = 0
	GAP      = 1
	MISMATCH = 1
)

type matrix [][]int

func newMatrix(n, m int) matrix {
	nm := make([][]int, n)
	for i := range nm {
		nm[i] = make([]int, m)
	}
	return nm
}

func score(seq1, seq2 string) int {
	n, m := len(seq1), len(seq2)
	nm := newMatrix(n+1, m+1)
	for i := range seq1 {
		nm[i][0] = GAP * i
	}
	for i := range seq2 {
		nm[0][i] = GAP * i
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			hor := nm[i-1][j] + GAP
			ver := nm[i][j-1] + GAP
			var diag int
			if seq1[i-1] == seq2[j-1] {
				diag = MATCH
			} else {
				diag = MISMATCH
			}
			diag += nm[i-1][j-1]
			nm[i][j] = min(hor, ver, diag)
		}
	}
	return nm[n][m]
}

type bkTree struct {
	term     string
	children map[int]*bkTree
	deleted  bool
}

func (b *bkTree) search(term string, fuzziness int) bool {
	d0 := score(b.term, term)
	if d0 <= fuzziness && !b.deleted {
		return true
	}
	low := d0 - fuzziness
	high := d0 + fuzziness
	for dist, node := range b.children {
		if dist >= low && dist <= high {
			if result := node.search(term, fuzziness); result {
				return true
			}
		}
	}
	return false
}

type BKTree struct {
	root      *bkTree
	Fuzziness int
}

func (b *BKTree) Add(term string) {
	q := &bkTree{
		term:     term,
		children: make(map[int]*bkTree),
	}
	if b.root == nil {
		b.root = q
		return
	}
	curr := b.root
	for {
		k := score(curr.term, term)
		if k == 0 {
			return
		}
		if child, ok := curr.children[k]; ok {
			curr = child
		} else {
			curr.children[k] = q
			return
		}
	}
}

func (b *BKTree) Has(term string) bool {
	return b.root.search(term, b.Fuzziness)
}

func (b *bkTree) deleteExact(term string) bool {
	k := score(b.term, term)
	if k == 0 {
		if b.deleted {
			return false
		}
		b.deleted = true
		return true
	}
	child, ok := b.children[k]
	if !ok {
		return false
	}
	return child.deleteExact(term)
}

func (b *BKTree) Del(term string) {
	if b.root != nil {
		b.root.deleteExact(term)
	}
}
